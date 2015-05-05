# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import datetime
import hashlib
import StringIO
import unittest
import zipfile

from google.appengine.ext import ndb
from testing_utils import testing

from components import auth
from components import utils

from cas import impl as cas_impl

from cipd import impl
from cipd import processing


class TestValidators(unittest.TestCase):
  def test_is_valid_package_name(self):
    self.assertTrue(impl.is_valid_package_name('a'))
    self.assertTrue(impl.is_valid_package_name('a/b'))
    self.assertTrue(impl.is_valid_package_name('a/b/c/1/2/3'))
    self.assertTrue(impl.is_valid_package_name('infra/tools/cipd'))
    self.assertTrue(impl.is_valid_package_name('-/_'))
    self.assertFalse(impl.is_valid_package_name(''))
    self.assertFalse(impl.is_valid_package_name('/a'))
    self.assertFalse(impl.is_valid_package_name('a/'))
    self.assertFalse(impl.is_valid_package_name('A'))
    self.assertFalse(impl.is_valid_package_name('a/B'))
    self.assertFalse(impl.is_valid_package_name('a\\b'))

  def test_is_valid_package_path(self):
    self.assertTrue(impl.is_valid_package_path('a'))
    self.assertTrue(impl.is_valid_package_path('a/b'))
    self.assertTrue(impl.is_valid_package_path('a/b/c/1/2/3'))
    self.assertTrue(impl.is_valid_package_path('infra/tools/cipd'))
    self.assertTrue(impl.is_valid_package_path('-/_'))
    self.assertFalse(impl.is_valid_package_path(''))
    self.assertFalse(impl.is_valid_package_path('/a'))
    self.assertFalse(impl.is_valid_package_path('a/'))
    self.assertFalse(impl.is_valid_package_path('A'))
    self.assertFalse(impl.is_valid_package_path('a/B'))
    self.assertFalse(impl.is_valid_package_path('a\\b'))

  def test_is_valid_instance_id(self):
    self.assertTrue(impl.is_valid_instance_id('a'*40))
    self.assertFalse(impl.is_valid_instance_id(''))
    self.assertFalse(impl.is_valid_instance_id('A'*40))

  def test_is_valid_instance_tag(self):
    self.assertTrue(impl.is_valid_instance_tag('k:v'))
    self.assertTrue(impl.is_valid_instance_tag('key:'))
    self.assertTrue(impl.is_valid_instance_tag('key-_01234:#$%@\//%$SD'))
    self.assertFalse(impl.is_valid_instance_tag(''))
    self.assertFalse(impl.is_valid_instance_tag('key'))
    self.assertFalse(impl.is_valid_instance_tag('KEY:'))
    self.assertFalse(impl.is_valid_instance_tag('key:' + 'a'*500))


class TestRepoService(testing.AppengineTestCase):
  maxDiff = None

  def setUp(self):
    super(TestRepoService, self).setUp()
    self.mocked_cas_service = MockedCASService()
    self.mock(impl.cas, 'get_cas_service', lambda: self.mocked_cas_service)
    self.service = impl.get_repo_service()

  def test_register_package_new(self):
    self.assertIsNone(self.service.get_package('a/b'))
    inst, registered = self.service.register_package(
        package_name='a/b',
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2014, 1, 1, 0, 0))
    self.assertTrue(registered)
    self.assertEqual('a/b', inst.package_name)
    self.assertEqual({
      'registered_by': auth.Identity(kind='user', name='abc@example.com'),
      'registered_ts': datetime.datetime(2014, 1, 1, 0, 0)
    }, inst.to_dict())

  def test_register_package_existing(self):
    inst, registered = self.service.register_package(
        package_name='a/b',
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2014, 1, 1, 0, 0))
    self.assertTrue(registered)
    inst, registered = self.service.register_package(
        package_name='a/b',
        caller=auth.Identity.from_bytes('user:def@example.com'),
        now=datetime.datetime(2015, 1, 1, 0, 0))
    self.assertFalse(registered)
    self.assertEqual({
      'registered_by': auth.Identity(kind='user', name='abc@example.com'),
      'registered_ts': datetime.datetime(2014, 1, 1, 0, 0)
    }, inst.to_dict())

  def test_register_instance_new(self):
    self.assertIsNone(self.service.get_instance('a/b', 'a'*40))
    self.assertIsNone(self.service.get_package('a/b'))
    inst, registered = self.service.register_instance(
        package_name='a/b',
        instance_id='a'*40,
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2014, 1, 1, 0, 0))
    self.assertTrue(registered)
    self.assertEqual(
        ndb.Key('Package', 'a/b', 'PackageInstance', 'a'*40), inst.key)
    self.assertEqual('a/b', inst.package_name)
    self.assertEqual('a'*40, inst.instance_id)
    expected = {
      'registered_by': auth.Identity(kind='user', name='abc@example.com'),
      'registered_ts': datetime.datetime(2014, 1, 1, 0, 0),
      'processors_failure': [],
      'processors_pending': [],
      'processors_success': [],
    }
    self.assertEqual(expected, inst.to_dict())
    self.assertEqual(
        expected, self.service.get_instance('a/b', 'a'*40).to_dict())
    self.assertTrue(self.service.get_package('a/b'))

  def test_register_instance_existing(self):
    # First register a package.
    inst1, registered = self.service.register_instance(
        package_name='a/b',
        instance_id='a'*40,
        caller=auth.Identity.from_bytes('user:abc@example.com'))
    self.assertTrue(registered)
    # Try to register it again.
    inst2, registered = self.service.register_instance(
          package_name='a/b',
          instance_id='a'*40,
          caller=auth.Identity.from_bytes('user:def@example.com'))
    self.assertFalse(registered)
    self.assertEqual(inst1.to_dict(), inst2.to_dict())

  def test_generate_fetch_url(self):
    inst, registered = self.service.register_instance(
        package_name='a/b',
        instance_id='a'*40,
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2014, 1, 1, 0, 0))
    self.assertTrue(registered)
    url = self.service.generate_fetch_url(inst)
    self.assertEqual(
        'https://signed-url/SHA1/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', url)

  def test_is_instance_file_uploaded(self):
    self.mocked_cas_service.uploaded[('SHA1', 'a'*40)] = ''
    self.assertTrue(self.service.is_instance_file_uploaded('a/b', 'a'*40))
    self.assertFalse(self.service.is_instance_file_uploaded('a/b', 'b'*40))

  def test_create_upload_session(self):
    upload_url, upload_session_id = self.service.create_upload_session(
        'a/b', 'a'*40, auth.Identity.from_bytes('user:abc@example.com'))
    self.assertEqual('http://upload_url', upload_url)
    self.assertEqual('upload_session_id', upload_session_id)

  def test_register_instance_with_processing(self):
    self.mock(utils, 'utcnow', lambda: datetime.datetime(2014, 1, 1))

    self.service.processors.append(MockedProcessor('bad', 'Error message'))
    self.service.processors.append(MockedProcessor('good'))

    tasks = []
    def mocked_enqueue_task(**kwargs):
      tasks.append(kwargs)
      return True
    self.mock(impl.utils, 'enqueue_task', mocked_enqueue_task)

    # The processors are added to the pending list.
    inst, registered = self.service.register_instance(
        package_name='a/b',
        instance_id='a'*40,
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2014, 1, 1, 0, 0))
    self.assertTrue(registered)
    expected = {
      'registered_by': auth.Identity(kind='user', name='abc@example.com'),
      'registered_ts': datetime.datetime(2014, 1, 1, 0, 0),
      'processors_failure': [],
      'processors_pending': ['bad', 'good'],
      'processors_success': [],
    }
    self.assertEqual(expected, inst.to_dict())

    # The processing task is enqueued.
    self.assertEqual([{
      'payload': '{"instance_id": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", '
          '"package_name": "a/b", "processors": ["bad", "good"]}',
      'queue_name': 'cipd-process',
      'transactional': True,
      'url': '/internal/taskqueue/cipd-process/'
          'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',
    }], tasks)

    # Now execute the task.
    self.service.process_instance(
        package_name='a/b',
        instance_id='a'*40,
        processors=['bad', 'good'])

    # Assert the final state.
    inst = self.service.get_instance('a/b', 'a'*40)
    expected = {
      'registered_by': auth.Identity(kind='user', name='abc@example.com'),
      'registered_ts': datetime.datetime(2014, 1, 1, 0, 0),
      'processors_failure': ['bad'],
      'processors_pending': [],
      'processors_success': ['good'],
    }
    self.assertEqual(expected, inst.to_dict())

    good_result = self.service.get_processing_result('a/b', 'a'*40, 'good')
    self.assertEqual({
      'created_ts': datetime.datetime(2014, 1, 1),
      'error': None,
      'result': {
        'instance_id': 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',
        'package_name': 'a/b',
        'processor_name': 'good',
      },
      'success': True,
    }, good_result.to_dict())

    bad_result = self.service.get_processing_result('a/b', 'a'*40, 'bad')
    self.assertEqual({
      'created_ts': datetime.datetime(2014, 1, 1),
      'error': 'Error message',
      'result': None,
      'success': False,
    }, bad_result.to_dict())

  def test_client_binary_extraction(self):
    self.mock(utils, 'utcnow', lambda: datetime.datetime(2014, 1, 1))

    # Prepare fake cipd binary package.
    out = StringIO.StringIO()
    zf = zipfile.ZipFile(out, 'w', zipfile.ZIP_DEFLATED)
    zf.writestr('cipd', 'cipd binary data here')
    zf.close()
    zipped = out.getvalue()
    digest = hashlib.sha1(zipped).hexdigest()

    # Pretend it is uploaded.
    self.mocked_cas_service.uploaded[('SHA1', digest)] = zipped

    # Register it as a package instance.
    self.mock(impl.utils, 'enqueue_task', lambda **_args: True)
    inst, registered = self.service.register_instance(
        package_name='infra/tools/cipd/linux-amd64',
        instance_id=digest,
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2014, 1, 1, 0, 0))
    self.assertTrue(registered)
    expected = {
      'registered_by': auth.Identity(kind='user', name='abc@example.com'),
      'registered_ts': datetime.datetime(2014, 1, 1, 0, 0),
      'processors_failure': [],
      'processors_pending': ['cipd_client_binary:v1'],
      'processors_success': [],
    }
    self.assertEqual(expected, inst.to_dict())

    # get_client_binary_info indicated that processing is not done yet.
    instance = self.service.get_instance(
        package_name='infra/tools/cipd/linux-amd64',
        instance_id=digest)
    info, error_msg = self.service.get_client_binary_info(instance)
    self.assertIsNone(info)
    self.assertIsNone(error_msg)

    # Execute post-processing task: it would extract CIPD binary.
    self.service.process_instance(
        package_name='infra/tools/cipd/linux-amd64',
        instance_id=digest,
        processors=['cipd_client_binary:v1'])

    # Ensure succeeded.
    result = self.service.get_processing_result(
        package_name='infra/tools/cipd/linux-amd64',
        instance_id=digest,
        processor_name='cipd_client_binary:v1')
    self.assertEqual({
      'created_ts': datetime.datetime(2014, 1, 1, 0, 0),
      'success': True,
      'error': None,
      'result': {
        'client_binary': {
          'hash_algo': 'SHA1',
          'hash_digest': '5a72c1535f8d132c341585207504d94e68ef8a9d',
          'size': 21,
        },
      },
    }, result.to_dict())

    # Verify get_client_binary_info works too.
    instance = self.service.get_instance(
        package_name='infra/tools/cipd/linux-amd64',
        instance_id=digest)
    info, error_msg = self.service.get_client_binary_info(instance)
    expected = impl.ClientBinaryInfo(
        sha1='5a72c1535f8d132c341585207504d94e68ef8a9d',
        size=21,
        fetch_url=(
            'https://signed-url/SHA1/5a72c1535f8d132c341585207504d94e68ef8a9d'))
    self.assertIsNone(error_msg)
    self.assertEqual(expected, info)

  def test_client_binary_extract_failure(self):
    self.mock(utils, 'utcnow', lambda: datetime.datetime(2014, 1, 1))

    # Pretend some fake data is uploaded.
    self.mocked_cas_service.uploaded[('SHA1', 'a'*40)] = 'not a zip'

    # Register it as a package instance.
    self.mock(impl.utils, 'enqueue_task', lambda **_args: True)
    inst, registered = self.service.register_instance(
        package_name='infra/tools/cipd/linux-amd64',
        instance_id='a'*40,
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2014, 1, 1, 0, 0))
    self.assertTrue(registered)
    expected = {
      'registered_by': auth.Identity(kind='user', name='abc@example.com'),
      'registered_ts': datetime.datetime(2014, 1, 1, 0, 0),
      'processors_failure': [],
      'processors_pending': ['cipd_client_binary:v1'],
      'processors_success': [],
    }
    self.assertEqual(expected, inst.to_dict())

    # Execute post-processing task: it would fail extracting CIPD binary.
    self.service.process_instance(
        package_name='infra/tools/cipd/linux-amd64',
        instance_id='a'*40,
        processors=['cipd_client_binary:v1'])

    # Ensure error is reported.
    result = self.service.get_processing_result(
        package_name='infra/tools/cipd/linux-amd64',
        instance_id='a'*40,
        processor_name='cipd_client_binary:v1')
    self.assertEqual({
      'created_ts': datetime.datetime(2014, 1, 1, 0, 0),
      'success': False,
      'error': 'File is not a zip file',
      'result': None,
    }, result.to_dict())

    # Verify get_client_binary_info reports it too.
    instance = self.service.get_instance(
        package_name='infra/tools/cipd/linux-amd64',
        instance_id='a'*40)
    info, error_msg = self.service.get_client_binary_info(instance)
    self.assertIsNone(info)
    self.assertEqual(
        'Failed to extract the binary: File is not a zip file', error_msg)

  def test_attach_detach_tags(self):
    _, registered = self.service.register_instance(
        package_name='a/b',
        instance_id='a'*40,
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2014, 1, 1, 0, 0))
    self.assertTrue(registered)

    # Add a tag.
    attached = self.service.attach_tags(
        package_name='a/b',
        instance_id='a'*40,
        tags=['tag1:value1'],
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2014, 1, 1, 0, 0))
    self.assertEqual(
      {
        'tag1:value1': {
          'registered_by': auth.Identity(kind='user', name='abc@example.com'),
          'registered_ts': datetime.datetime(2014, 1, 1, 0, 0),
          'tag': 'tag1:value1',
        },
      }, {k: e.to_dict() for k, e in attached.iteritems()})
    self.assertEqual('a/b', attached['tag1:value1'].package_name)
    self.assertEqual('a'*40, attached['tag1:value1'].instance_id)

    # Attempt to attach existing one (and one new).
    attached = self.service.attach_tags(
        package_name='a/b',
        instance_id='a'*40,
        tags=['tag1:value1', 'tag2:value2'],
        caller=auth.Identity.from_bytes('user:abc@example.com'),
        now=datetime.datetime(2015, 1, 1, 0, 0))
    self.assertEqual(
      {
        'tag1:value1': {
          'registered_by': auth.Identity(kind='user', name='abc@example.com'),
          # Didn't change to 2015.
          'registered_ts': datetime.datetime(2014, 1, 1, 0, 0),
          'tag': 'tag1:value1',
        },
        'tag2:value2': {
          'registered_by': auth.Identity(kind='user', name='abc@example.com'),
          'registered_ts': datetime.datetime(2015, 1, 1, 0, 0),
          'tag': 'tag2:value2',
        },
      }, {k: e.to_dict() for k, e in attached.iteritems()})

    # Get specific tags.
    tags = self.service.get_tags('a/b', 'a'*40, ['tag1:value1', 'missing:'])
    self.assertEqual(
      {
        'tag1:value1': {
          'registered_by': auth.Identity(kind='user', name='abc@example.com'),
          'registered_ts': datetime.datetime(2014, 1, 1, 0, 0),
          'tag': 'tag1:value1',
        },
        'missing:': None,
      }, {k: e.to_dict() if e else None for k, e in tags.iteritems()})

    # Get all tags. Newest first.
    tags = self.service.query_tags('a/b', 'a'*40)
    self.assertEqual(['tag2:value2', 'tag1:value1'], [t.tag for t in tags])

    # Search by specific tag (in a package).
    found = self.service.search_by_tag('tag1:value1', package_name='a/b')
    self.assertEqual(
        [('a/b', 'a'*40)], [(e.package_name, e.instance_id) for e in found])

    # Search by specific tag (globally). Use callback to cover this code path.
    found = self.service.search_by_tag('tag1:value1')
    self.assertEqual(
        [('a/b', 'a'*40)], [(e.package_name, e.instance_id) for e in found])

    # Cover callback usage.
    found = self.service.search_by_tag(
        'tag1:value1', callback=lambda *_a: False)
    self.assertFalse(found)

    # Remove tag, search again -> missing.
    self.service.detach_tags('a/b', 'a'*40, ['tag1:value1', 'missing:'])
    found = self.service.search_by_tag('tag1:value1')
    self.assertFalse(found)


class MockedCASService(object):
  def __init__(self):
    self.uploaded = {}

  def is_fetch_configured(self):
    return True

  def generate_fetch_url(self, algo, digest):
    return 'https://signed-url/%s/%s' % (algo, digest)

  def is_object_present(self, algo, digest):
    return (algo, digest) in self.uploaded

  def create_upload_session(self, _algo, _digest, _caller):
    class UploadSession(object):
      upload_url = 'http://upload_url'
    return UploadSession(), 'upload_session_id'

  def open(self, hash_algo, hash_digest, read_buffer_size):
    assert read_buffer_size > 0
    if not self.is_object_present(hash_algo, hash_digest):  # pragma: no cover
      raise cas_impl.NotFoundError()
    return StringIO.StringIO(self.uploaded[(hash_algo, hash_digest)])

  def start_direct_upload(self, hash_algo):
    assert hash_algo == 'SHA1'
    return cas_impl.DirectUpload(
        file_obj=StringIO.StringIO(),
        hasher=hashlib.sha1(),
        callback=lambda *_args: None)


class MockedProcessor(processing.Processor):
  def __init__(self, name, error=None):
    self.name = name
    self.error = error

  def should_process(self, instance):
    return True

  def run(self, instance, data):
    if self.error:
      raise processing.ProcessingError(self.error)
    return {
      'instance_id': instance.instance_id,
      'package_name': instance.package_name,
      'processor_name': self.name,
    }
