# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import collections
import logging
import posixpath
import sys

from infra.libs.git2 import INVALID
from infra.libs.git2 import config_ref
from infra.libs.git2 import repo


LOGGER = logging.getLogger(__name__)


################################################################################
# ConfigRef
################################################################################

class GsubtreedConfigRef(config_ref.ConfigRef):
  CONVERT = {
    'interval': lambda self, val: float(val),
    'subtree_synthesized_prefix': lambda self, val: str(val),
    'subtree_processed_prefix': lambda self, val: str(val),

    'base_url': lambda self, val: str(val) if val else self.repo.url,
    'enabled_refglobs': lambda self, val: map(str, list(val)),
    # normpath to avoid trailing/double-slash errors.
    'enabled_paths': lambda self, val: map(posixpath.normpath, map(str, val)),
  }
  DEFAULTS = {
    'interval': 5.0,

    # e.g. while processing the subtree 'b/foo' on refs/heads/master
    #   refs/heads/master                              <- real commits
    #   refs/subtree-processed/b/foo/-/heads/master    <- ancestor tag of master
    #   refs/subtree-synthesized/b/foo/-/heads/master  <- ref with synth commits
    # For the sake of implementation simplicity, this daemon assumes the
    # googlesource.com guarantee of transactional multi-ref pushes within a
    # single repo.
    'subtree_processed_prefix': 'refs/subtree-processed',
    'subtree_synthesized_prefix': 'refs/subtree-synthesized',

    # The base URL is the url relative to which all mirror repos are assumed to
    # exist. For example, if you mirror the path 'bob', and base_url is
    # https://host.domain.tld/main_repo, then it would assume that the mirror
    # for the bob subtree is https://host.domain.tld/main_repo/bob.
    #
    # By default, base_url is set to the repo that gsubtreed is processing
    'base_url': None,
    'enabled_refglobs': ['refs/heads/*'],
    'enabled_paths': [],
  }
  REF = 'refs/gsubtreed-config/main'



################################################################################
# Core functionality
################################################################################

def process_path(path, origin_repo, config):
  def join(prefix, ref):
    assert ref.ref.startswith('refs/')
    ref = '/'.join((prefix, path)) + '/-/' + ref.ref[len('refs/'):]
    return origin_repo[ref]

  origin_push = {}

  base_url = config['base_url']
  mirror_url = '[FILE-URL]' if base_url.startswith('file:') else origin_repo.url

  subtree_repo = repo.Repo(posixpath.join(base_url, path))
  subtree_repo.repos_dir = origin_repo.repos_dir
  subtree_repo.reify(share_from=origin_repo)
  subtree_repo.run('fetch', stdout=sys.stdout, stderr=sys.stderr)
  subtree_repo_push = {}

  synthed_count = 0

  for glob in config['enabled_refglobs']:
    for ref in origin_repo.refglob(glob):
      LOGGER.info('processing ref %s', ref)
      processed = join(config['subtree_processed_prefix'], ref)
      synthed = join(config['subtree_synthesized_prefix'], ref)

      synth_parent = synthed.commit
      LOGGER.info('starting with tree %r', synthed.commit.data.tree)

      for commit in processed.to(ref, path):
        LOGGER.info('processing commit %s', commit)
        obj_name = '{.hsh}:{}'.format(commit, path)
        typ = origin_repo.run('cat-file', '-t', obj_name).strip()
        if typ != 'tree':
          LOGGER.warn('path %r is not a tree in commit %s', path, commit)
          continue
        dir_tree = origin_repo.run('rev-parse', obj_name).strip()

        LOGGER.info('found new tree %r', dir_tree)
        synthed_count += 1
        synth_parent = commit.alter(
          parents=[synth_parent.hsh] if synth_parent is not INVALID else [],
          tree=dir_tree,
          footers=collections.OrderedDict([
            ('Cr-Mirrored-From', [mirror_url]),
            ('Cr-Mirrored-Commit', [commit.hsh]),
          ]),
        )
        origin_push[synthed] = synth_parent
        subtree_repo_push[subtree_repo[ref.ref]] = synth_parent

      origin_push[processed] = ref.commit

  success = True
  # TODO(iannucci): Return the pushspecs from this method, and then thread
  # the dispatches to subtree_repo. Additionally, can batch the origin_repo
  # pushes (and push them serially in batches as the subtree_repo pushes
  # complete).

  # because the hashes are deterministic based on the real history, the pushes
  # can happen completely independently. If we miss one, we'll catch it on the
  # next pass.
  try:
    origin_repo.fast_forward_push(origin_push)
  except Exception:  # pragma: no cover
    LOGGER.exception('Caught exception while pushing origin in process_path')
    success = False

  try:
    subtree_repo.fast_forward_push(subtree_repo_push)
  except Exception:  # pragma: no cover
    LOGGER.exception('Caught exception while pushing subtree in process_path')
    success = False

  return success, synthed_count


def inner_loop(origin_repo, config):
  """Returns (success, {path: #commits_synthesized})."""

  LOGGER.debug('fetching %r', origin_repo)
  origin_repo.run('fetch', stdout=sys.stdout, stderr=sys.stderr)
  config.evaluate()

  success = True
  processed = {}
  for path in config['enabled_paths']:
    LOGGER.info('processing path %s', path)
    try:
      path_success, num_synthed = process_path(path, origin_repo, config)
      success = path_success and success
      processed[path] = num_synthed
    except Exception:  # pragma: no cover
      LOGGER.exception('Caught in inner_loop')
      success = False

  return success, processed
