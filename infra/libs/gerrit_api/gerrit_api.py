# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

"""Yet another wrapper around Gerrit REST API."""

import base64
import cookielib
import json
import logging
import netrc
import os
import requests
import requests_cache
import stat
import sys
import time
import urllib

from requests.packages import urllib3


LOGGER = logging.getLogger(__name__)


class GerritException(Exception):
  """Base class for exceptions raised by this module."""


class NetrcException(GerritException):
  """Netrc file is missing or incorrect."""


class GitcookiesException(GerritException):
  """Gitcookies file is missing or incorrect."""


class UnexpectedResponseException(GerritException):
  """Gerrit returned something unexpected."""

  def __init__(self, http_code, body):  # pragma: no cover
    super(UnexpectedResponseException, self).__init__()
    self.http_code = http_code
    self.body = body

  def __str__(self):  # pragma: no cover
    return 'Unexpected response (HTTP %d): %s' % (self.http_code, self.body)


class BlockCookiesPolicy(cookielib.DefaultCookiePolicy):
  def set_ok(self, cookie, request):
    return False


class Gerrit(object):  # pragma: no cover
  """Wrapper around a single Gerrit host.

  Args:
    host (str): gerrit host name.
    netrc_path (str): path to local netrc file. If None, the default location
      for the current OS is used.
    gitcookies_path (str): path to local gitcookies file. If None, the default
      location for the current OS is used.
    throttle_delay_sec (int): minimal time delay between two requests, to
      avoid hammering the Gerrit server.
  """

  def __init__(self, host, netrc_path=None, gitcookies_path=None,
               throttle_delay_sec=0):
    auth = _load_netrc(netrc_path).authenticators(host)
    if not auth:
      auth = _load_gitcookies(gitcookies_path, host)
    if not auth:
      raise GerritException('No auth record for %s' % host)
    self._auth_header = 'Basic %s' % (
        base64.b64encode('%s:%s' % (auth[0], auth[2])))
    self._url_base = 'https://%s/a' % host.rstrip('/')
    self._throttle = throttle_delay_sec
    self._last_call_ts = None
    self.session = requests.Session()
    # Do not use cookies with Gerrit. This breaks interaction with Google's
    # Gerrit instances. Do not use cookies as advised by the Gerrit team.
    self.session.cookies.set_policy(BlockCookiesPolicy())
    retry_config = urllib3.util.Retry(total=4, backoff_factor=0.5,
                                      status_forcelist=[500, 503])
    self.session.mount(self._url_base, requests.adapters.HTTPAdapter(
        max_retries=retry_config))


  def _request(self, method, request_path, params=None, body=None):
    """Sends HTTP request to Gerrit.

    Args:
      method: HTTP method (e.g 'GET', 'POST', ...).
      request_path: URL of the endpoint, relative to host (e.g. '/accounts/id').
      params: dict with query parameters.
      body: optional request body, will be serialized to JSON.

    Returns:
      Tuple (response code, deserialized JSON response).
    """
    if not request_path.startswith('/'):
      request_path = '/' + request_path

    full_url = '%s%s' % (self._url_base, request_path)

    # Wait to avoid Gerrit quota, don't wait if a response is in the cache.
    if self._throttle and not _is_response_cached(method, full_url):
      now = time.time()
      if self._last_call_ts and now - self._last_call_ts < self._throttle:
        time.sleep(self._throttle - (now - self._last_call_ts))
      self._last_call_ts = time.time()

    headers = {
        # This makes the server return compact JSON.
        'Accept': 'application/json',
        # This means responses will be gzip compressed.
        'Accept-encoding': 'gzip',
        'Authorization': self._auth_header,
    }

    if body is not None:
      body = json.dumps(body)
      headers['Content-Type'] = 'application/json;charset=UTF-8'
      headers['Content-Length'] = str(len(body))

    LOGGER.debug('%s %s', method, full_url)
    response = self.session.request(
        method=method,
        url=full_url,
        params=params,
        data=body,
        headers=headers)

    # Gerrit prepends )]}' to response.
    prefix = ')]}\'\n'
    body = response.text
    if body and body.startswith(prefix):
      body = json.loads(body[len(prefix):])

    return response.status_code, body

  def get_account(self, account_id):
    """Returns a dict describing a Gerrit account or None if no such account.

    Args:
      account_id: email, numeric account id, or 'self'.

    Returns:
      None if no such account, AccountInfo dict otherwise.
    """
    assert '/' not in account_id
    code, body = self._request('GET', '/accounts/%s' % account_id)
    if code == 200:
      return body
    if code == 404:
      return None
    raise UnexpectedResponseException(code, body)

  def add_group_members(self, group, members):
    """Adds a bunch of members to a group.

    Args:
      group: name of a group to add members to.
      members: iterable with emails of accounts to add to the group.

    Returns:
      None

    Raises:
      UnexpectedResponseException: if call failed.
    """
    if '/' in group:
      raise ValueError('Invalid group name: %s' % group)
    code, body = self._request(
        method='POST',
        request_path='/groups/%s/members.add' % group,
        body={'members': list(members)})
    if code != 200:
      raise UnexpectedResponseException(code, body)

  def is_account_active(self, account_id):
    if '/' in account_id:
      raise ValueError('Invalid account id: %s' % account_id)
    code, body = self._request(
        method='GET',
        request_path='/accounts/%s/active' % account_id)
    if code == 200:
      return True
    if code == 204:
      return False
    raise UnexpectedResponseException(code, body)

  def activate_account(self, account_id):
    """Sets account state to 'active'.

    Args:
      account_id (str): account to update.

    Raises:
      UnexpectedResponseException: if gerrit does not answer as expected.
    """
    if '/' in account_id:
      raise ValueError('Invalid account id: %s' % account_id)
    code, body = self._request(
        method='PUT',
        request_path='/accounts/%s/active' % account_id)
    if code not in (200, 201):
      raise UnexpectedResponseException(code, body)

  def get_projects(self, prefix=''):
    """Returns list of projects with names starting with a prefix.

    Args:
      prefix (str): optional project name prefix to limit the listing to.

    Returns:
      Dict <project name> -> {'state': 'ACTIVE', 'parent': 'All-Projects'}

    Raises:
      UnexpectedResponseException: if gerrit does not answer as expected.
    """
    code, body = self._request(
        method='GET',
        request_path='/projects/?p=%s&t' % urllib.quote(prefix, safe=''))
    if code not in (200, 201):
      raise UnexpectedResponseException(code, body)
    return body

  def get_project_parent(self, project):
    """Retrieves the name of a project's parent project.

    Returns None If |project| is not registered on Gerrit or doesn't have
    a parent (like 'All-Projects').

    Args:
      project (str): project to query.

    Raises:
      UnexpectedResponseException: if gerrit does not answer as expected.
    """
    code, body = self._request(
        method='GET',
        request_path='/projects/%s/parent' % urllib.quote(project, safe=''))
    if code == 404:
      return None
    if code not in (200, 201):
      raise UnexpectedResponseException(code, body)
    assert isinstance(body, unicode)
    return body if body else None

  def set_project_parent(self, project, parent, commit_message=None):
    """Changes project's parent project.

    Args:
      project (str): project to change.
      parent (str): parent to set.
      commit_message (str): message for corresponding refs/meta/config commit.

    Raises:
      UnexpectedResponseException: if gerrit does not answer as expected.
    """
    commit_message = (
        commit_message or ('Changing parent project to %s' % parent))
    code, body = self._request(
        method='PUT',
        request_path='/projects/%s/parent' % urllib.quote(project, safe=''),
        body={'parent': parent, 'commit_message': commit_message})
    if code not in (200, 201):
      raise UnexpectedResponseException(code, body)


def _load_netrc(path=None):  # pragma: no cover
  """Loads netrc file with gerrit credentials.

  Args:
    path: path to .netrc or None to use default path.

  Returns:
    netrc_obj (:class:`netrc.netrc`):

  Raises:
    NetrcException: if the netrc file can't be read, for any reason.
  """
  if not path:
    # HOME might not be set on Windows.
    if 'HOME' not in os.environ:
      raise NetrcException('HOME environment variable is not set')
    path = os.path.join(
        os.environ['HOME'],
        '_netrc' if sys.platform.startswith('win') else '.netrc')
  try:
    return netrc.netrc(path)
  except IOError as exc:
    raise NetrcException('Could not read netrc file %s: %s' % (path, exc))
  except netrc.NetrcParseError as exc:
    netrc_stat = os.stat(exc.filename)
    if netrc_stat.st_mode & (stat.S_IRWXG | stat.S_IRWXO):
      raise NetrcException(
          'netrc file %s cannot be used because its file permissions '
          'are insecure.  netrc file permissions should be 600.' % path)
    else:
      raise NetrcException(
          'Cannot use netrc file %s due to a parsing error: %s' % (path, exc))


def _load_gitcookies(path, host):
  """Loads gitcookies file with gerrit credentials.

  Args:
    path: path to .gitcookies or None to use default path.

  Returns:
    returns (login, None, password) tuple representing the credentials.
  """
  if not path:
    # HOME might not be set on Windows.
    if 'HOME' not in os.environ:
      raise GitcookiesException('HOME environment variable is not set')
    path = os.path.join(os.environ['HOME'], '.gitcookies')

  try:
    with open(path) as f:
      for line in f:
        fields = line.strip().split('\t')
        if line.strip().startswith('#') or len(fields) != 7:
          continue
        domain, xpath, key, value = fields[0], fields[2], fields[5], fields[6]
        if cookielib.domain_match(host, domain) and xpath == '/' and key == 'o':
          login, password = value.split('=', 1)
          return (login, None, password)
  except IOError:
    pass
  return None


def _is_response_cached(method, full_url):  # pragma: no cover
  """Returns True if response to GET request is in requests_cache.

  Args:
    method (str): http verb ('GET', 'POST', etc.)
    full_url (str): url, including the protocol
  Returns:
    is_cached (bool):
"""
  if method != 'GET':
    return False
  try:
    cache = requests_cache.get_cache()
  except AttributeError:
    cache = None
  return cache.has_url(full_url) if cache else False

