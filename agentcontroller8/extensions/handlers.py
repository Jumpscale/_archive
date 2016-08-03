import json
import requests
import logging
import hashlib

import utils
from JumpScale import j

ENDPOINT_CONFIG = '/rest/system/config'
ENDPOINT_RESTART = '/rest/system/restart'


"""
SHARE_FOLDERS defines what folders to share on agent
key: local folder id
value: remote path to where to share on agent side.

Currently we only share legacy and jumpscripts folder to agent

SHARE_FOLDERS {
    folder_id: folder_path
}

folder_path is relative to agent-home but can be set as full path if started with /
"""
SHARE_FOLDERS = {
    # 'legacy': 'legacy',
    'jumpscripts': 'jumpscripts'
}

settings = {
    'syncthing': {
        'url': 'http://localhost:8384/',
    },
    'redis': {
        'address': 'localhost',
        'port': '6379',
        'password': None
    }
}


def init(config):
    """
    init gets called immediately when agent controller starts with the settings from the config file.
    We use this to configure the module for the operations needed.
    """
    settings['syncthing']['url'] = config['syncthing_url']
    settings['redis']['address'] = config['redis_address']
    settings['redis']['port'] = config['redis_port']
    settings['redis']['password'] = config['redis_password']


def results_or_die(results):
    if results.state != 'SUCCESS':
        raise Exception('Error executing cmd %s.%s: %s' % (results.cmd, results.args.name, results.data))
    assert results.level == 20, 'Only json response is supported so far'

    return json.loads(results.data)


def get_url(endpoint):
    base_url = settings['syncthing']['url'].rstrip('/')
    return '%s%s' % (base_url, endpoint)


def openPortForward(client, gid, nid):
    tunnels = client.tunnel_list(gid, nid)
    synctunnel = list(filter(lambda t: t['remote'] == 22000 and t['gateway'] == 'controller', tunnels))
    if synctunnel:
        tunnel = synctunnel[0]
    else:
        tunnel = client.tunnel_open(gid, nid, 0, 'controller', '127.0.0.1', 22000)

    return '127.0.0.1:%d' % tunnel['local']


@utils.exclusive('/tmp/agent-start.lock')
def startup(gid, nid):
    # TODO: client must use settings of somekind
    sessions = requests.Session()

    headers = {
        'content-type': 'application/json'
    }

    client = j.clients.agentcontroller.getAdvanced(**settings['redis'])
    default = j.clients.agentcontroller.getRunArgs(domain='jumpscale')

    get_id = client.cmd(gid, nid, 'sync', default.update({'name': 'get_id'}))

    address = openPortForward(client, gid, nid)
    agent_device_id = results_or_die(get_id.get_next_result(30))

    endpoint = get_url(ENDPOINT_CONFIG)
    response = sessions.get(endpoint, headers=headers)

    if not response.ok:
        raise Exception('Invalid response from syncthing', response.reason)

    local_device_id = response.headers['x-syncthing-id']
    device_id_hash = hashlib.md5(local_device_id.encode()).hexdigest()
    config = response.json()

    # if auth is off, we still need to use the API key to be able to use POST.
    # in this case, get the API key from the get response
    api_key = config['gui']['apiKey']
    headers['X-API-Key'] = api_key

    devices = list(filter(lambda d: d['deviceID'] == agent_device_id, config['devices']))

    dirty = False
    if not devices:
        device = {
            'addresses': ['dynamic'],
            'certName': '',
            'compression': 'metadata',
            'deviceID': agent_device_id,
            'introducer': False,
            'name': '%s-%s' % (gid, nid)
        }

        config['devices'].append(device)
        dirty = True

    # add device to shared folder.
    for folder_id_prefix in SHARE_FOLDERS:
        folder_id = '%s-%s' % (folder_id_prefix, device_id_hash)
        folders = list(filter(lambda f: f['id'] == folder_id, config['folders']))

        if not folders:
            logging.warn('Folder id "%s" is not shared on syncthing', folder_id)
            continue

        folder = folders[0]
        if not list(filter(lambda d: d['deviceID'] == agent_device_id, folder['devices'])):
            # share folder with device.

            folder['devices'].append({
                'deviceID': agent_device_id
            })
            dirty = True

    if dirty:
        # apply changes
        response = sessions.post(endpoint, data=json.dumps(config), headers=headers)
        if not response.ok:
            raise Exception('Failed to set syncthing configuration', response.reason)

        response = sessions.post(get_url(ENDPOINT_RESTART), headers=headers)
        if not response.ok:
            raise Exception('Failed to restart syncthing', get_url(ENDPOINT_RESTART), response.reason)

    # Now, the syncthing on AC side knows about the syncthing on Agent side. Now we have
    # to register this instance of syncthing on agent as well. We can do this via a simple agent command

    for folder_id_prefix, remote_path in SHARE_FOLDERS.items():
        # NOTE: the address is set to 127.0.0.1:33000 because the agent automatically opens a tunnel
        # to the master node (this machine)
        folder_id = '%s-%s' % (folder_id_prefix, device_id_hash)
        data = {
            'device_id': local_device_id,
            'folder_id': folder_id,
            'path': remote_path,
            'address': address,
        }

        result = client.cmd(gid, nid,
                            'sync', default.update({'name': 'sync_folder'}),
                            json.dumps(data)).get_next_result()

        if result.state != 'SUCCESS':
            logging.warn('Error syncthing jumpscripts folder on agent: %s' % result.data)
            continue

    client.cmd(gid, nid, 'sync', default.update({'name': 'restart'}))
