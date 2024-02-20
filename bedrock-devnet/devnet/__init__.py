import argparse
import logging
import os
import subprocess
import json
import socket
import calendar
import datetime
import time
import requests
import shutil
import http.client
from multiprocessing import Process, Queue
import concurrent.futures
from collections import namedtuple


import devnet.log_setup
from dotenv import dotenv_values

pjoin = os.path.join

parser = argparse.ArgumentParser(description='Bedrock devnet launcher')
parser.add_argument('--monorepo-dir', help='Directory of the monorepo', default=os.getcwd())
parser.add_argument('--allocs', help='Only create the allocs and exit', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument('--test', help='Tests the deployment, must already be deployed', type=bool, action=argparse.BooleanOptionalAction)

log = logging.getLogger()

class Bunch:
    def __init__(self, **kwds):
        self.__dict__.update(kwds)

class ChildProcess:
    def __init__(self, func, *args):
        self.errq = Queue()
        self.process = Process(target=self._func, args=(func, args))

    def _func(self, func, args):
        try:
            func(*args)
        except Exception as e:
            self.errq.put(str(e))

    def start(self):
        self.process.start()

    def join(self):
        self.process.join()

    def get_error(self):
        return self.errq.get() if not self.errq.empty() else None


def main():
    args = parser.parse_args()

    monorepo_dir = os.path.abspath(args.monorepo_dir)
    devnet_dir = pjoin(monorepo_dir, '.devnet')
    contracts_bedrock_dir = pjoin(monorepo_dir, 'packages', 'contracts-bedrock')
    deployment_dir = pjoin(contracts_bedrock_dir, 'deployments', 'devnetL1')
    op_node_dir = pjoin(args.monorepo_dir, 'op-node')
    ops_bedrock_dir = pjoin(monorepo_dir, 'ops-bedrock')
    deploy_config_dir = pjoin(contracts_bedrock_dir, 'deploy-config')
    devnet_config_path = pjoin(deploy_config_dir, 'devnetL1.json')
    devnet_config_template_path = pjoin(deploy_config_dir, 'devnetL1-template.json')
    ops_chain_ops = pjoin(monorepo_dir, 'op-chain-ops')
    sdk_dir = pjoin(monorepo_dir, 'packages', 'sdk')

    paths = Bunch(
      mono_repo_dir=monorepo_dir,
      devnet_dir=devnet_dir,
      contracts_bedrock_dir=contracts_bedrock_dir,
      deployment_dir=deployment_dir,
      l1_deployments_path=pjoin(deployment_dir, '.deploy'),
      deploy_config_dir=deploy_config_dir,
      devnet_config_path=devnet_config_path,
      devnet_config_template_path=devnet_config_template_path,
      op_node_dir=op_node_dir,
      ops_bedrock_dir=ops_bedrock_dir,
      ops_chain_ops=ops_chain_ops,
      sdk_dir=sdk_dir,
      genesis_l1_path=pjoin(devnet_dir, 'genesis-l1.json'),
      genesis_l2_path=pjoin(devnet_dir, 'genesis-l2.json'),
      allocs_path=pjoin(devnet_dir, 'allocs-l1.json'),
      addresses_json_path=pjoin(devnet_dir, 'addresses.json'),
      sdk_addresses_json_path=pjoin(devnet_dir, 'sdk-addresses.json'),
      rollup_config_path=pjoin(devnet_dir, 'rollup.json')
    )

    if args.test:
      log.info('Testing deployed devnet')
      devnet_test(paths)
      return

    os.makedirs(devnet_dir, exist_ok=True)

    if args.allocs:
        devnet_l1_genesis(paths)
        return

    git_commit = subprocess.run(['git', 'rev-parse', 'HEAD'], capture_output=True, text=True).stdout.strip()
    git_date = subprocess.run(['git', 'show', '-s', "--format=%ct"], capture_output=True, text=True).stdout.strip()

    # CI loads the images from workspace, and does not otherwise know the images are good as-is
#     if os.getenv('DEVNET_NO_BUILD') == "true":
#         log.info('Skipping docker images build')
#     else:
#         log.info(f'Building docker images for git commit {git_commit} ({git_date})')
#         run_command(['docker', 'compose', 'build', '--progress', 'plain',
#                      '--build-arg', f'GIT_COMMIT={git_commit}', '--build-arg', f'GIT_DATE={git_date}'],
#                     cwd=paths.ops_bedrock_dir, env={
#             'PWD': paths.ops_bedrock_dir,
#             'DOCKER_BUILDKIT': '1', # (should be available by default in later versions, but explicitly enable it anyway)
#             'COMPOSE_DOCKER_CLI_BUILD': '1'  # use the docker cache
#         })

    log.info('Devnet starting')
    devnet_deploy(paths)


def deploy_contracts(paths):
    wait_up(8545)
    wait_for_rpc_server('http://127.0.0.1:8545')
#     log.info('Wait for L1 for a period of time to avoid submitting transactions in the first few block heights.')
#     time.sleep(10)

    l1env = dotenv_values('./ops-bedrock/l1.env')
    log.info(l1env)
    l1_init_holder = l1env['INIT_HOLDER']
    l1_init_holder_prv = l1env['INIT_HOLDER_PRV']
    proposer_address = l1env['PROPOSER_ADDRESS']
    account = l1_init_holder
    log.info(f'Deploying with {account}')

    # send some ether to the create2 deployer account
    run_command([
        'cast', 'send', '--private-key', l1_init_holder_prv,
        '--rpc-url', 'http://127.0.0.1:8545', '--gas-price', '10000000000', '--legacy',
        '--value', '1ether', '0x3fAB184622Dc19b6109349B94811493BF2a45362'
    ], env={}, cwd=paths.contracts_bedrock_dir)

    # send some ether to proposer address
    run_command([
        'cast', 'send', '--private-key', l1_init_holder_prv,
        '--rpc-url', 'http://127.0.0.1:8545', '--gas-price', '10000000000', '--legacy',
        '--value', '10000ether', proposer_address
    ], env={}, cwd=paths.contracts_bedrock_dir)

    # deploy the create2 deployer
    run_command([
      'cast', 'publish', '--rpc-url', 'http://127.0.0.1:8545',
      '0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222'
    ], env={}, cwd=paths.contracts_bedrock_dir)

    fqn = 'scripts/Deploy.s.sol:Deploy'
    run_command([
        'forge', 'script', fqn, '--private-key', l1_init_holder_prv, '--with-gas-price', '10000000000', '--legacy',
        '--rpc-url', 'http://127.0.0.1:8545', '--broadcast',
    ], env={}, cwd=paths.contracts_bedrock_dir)

    shutil.copy(paths.l1_deployments_path, paths.addresses_json_path)

    log.info('Syncing contracts.')
    run_command([
        'forge', 'script', fqn, '--sig', 'sync()',
        '--rpc-url', 'http://127.0.0.1:8545'
    ], env={}, cwd=paths.contracts_bedrock_dir)

def init_devnet_l1_deploy_config(paths, update_timestamp=False):
    deploy_config = read_json(paths.devnet_config_template_path)
    if update_timestamp:
        deploy_config['l1GenesisBlockTimestamp'] = '{:#x}'.format(int(time.time()))
    write_json(paths.devnet_config_path, deploy_config)

def devnet_l1_genesis(paths):
    log.info('Starting L1.')
    init_devnet_l1_deploy_config(paths)

    run_command(['docker-compose', 'up', '-d', 'l1'], cwd=paths.ops_bedrock_dir, env={
        'PWD': paths.ops_bedrock_dir
    })

    forge = ChildProcess(deploy_contracts, paths)
    forge.start()
    forge.join()
    err = forge.get_error()
    if err:
        raise Exception(f"Exception occurred in child process: {err}")

    log.info('Start and Deploy L1 success.')


# Bring up the devnet where the contracts are deployed to L1
def devnet_deploy(paths):
    if os.path.exists(paths.addresses_json_path):
        log.info('L1 genesis already generated.')
        log.info('Starting L1.')
        init_devnet_l1_deploy_config(paths)

        run_command(['docker-compose', 'up', '-d', 'l1'], cwd=paths.ops_bedrock_dir, env={
            'PWD': paths.ops_bedrock_dir
        })
        wait_up(8545)
        wait_for_rpc_server('http://127.0.0.1:8545')
    else:
        log.info('Generating L1 genesis.')
        if os.path.exists(paths.allocs_path) == False:
            devnet_l1_genesis(paths)

    l1env = dotenv_values('./ops-bedrock/l1.env')
    log.info(l1env)
    bscChainId = l1env['BSC_CHAIN_ID']
    l1_init_holder = l1env['INIT_HOLDER']
    l1_init_holder_prv = l1env['INIT_HOLDER_PRV']
    proposer_address_prv = l1env['PROPOSER_ADDRESS_PRV']
    log.info('Generating network config.')
    devnet_cfg_orig = pjoin(paths.contracts_bedrock_dir, 'deploy-config', 'devnetL1.json')
    devnet_cfg_backup = pjoin(paths.devnet_dir, 'devnetL1.json.bak')
    shutil.copy(devnet_cfg_orig, devnet_cfg_backup)
    deploy_config = read_json(devnet_cfg_orig)
    l1BlockTag = l1BlockTagGet()["result"]
    log.info(l1BlockTag)
    l1BlockTimestamp = l1BlockTimestampGet(l1BlockTag)["result"]["timestamp"]
    log.info(l1BlockTimestamp)
    deploy_config['l1GenesisBlockTimestamp'] = l1BlockTimestamp
    deploy_config['l1StartingBlockTag'] = l1BlockTag
    deploy_config['l1ChainID'] = int(bscChainId,10)
    deploy_config['batchSenderAddress'] = l1_init_holder
    deploy_config['l2OutputOracleProposer'] = l1_init_holder
    deploy_config['baseFeeVaultRecipient'] = l1_init_holder
    deploy_config['l1FeeVaultRecipient'] = l1_init_holder
    deploy_config['sequencerFeeVaultRecipient'] = l1_init_holder
    deploy_config['proxyAdminOwner'] = l1_init_holder
    deploy_config['finalSystemOwner'] = l1_init_holder
    deploy_config['portalGuardian'] = l1_init_holder
    deploy_config['governanceTokenOwner'] = l1_init_holder
    write_json(devnet_cfg_orig, deploy_config)

    if os.path.exists(paths.genesis_l2_path):
        log.info('L2 genesis and rollup configs already generated.')
    else:
        log.info('Generating network config.')
        log.info('Generating L2 genesis and rollup configs.')
        run_command([
            'go', 'run', 'cmd/main.go', 'genesis', 'l2',
            '--l1-rpc', 'http://localhost:8545',
            '--deploy-config', paths.devnet_config_path,
            '--deployment-dir', paths.deployment_dir,
            '--outfile.l2', pjoin(paths.devnet_dir, 'genesis-l2.json'),
            '--outfile.rollup', pjoin(paths.devnet_dir, 'rollup.json')
        ], cwd=paths.op_node_dir)

    if os.path.exists(devnet_cfg_backup):
        shutil.move(devnet_cfg_backup, devnet_cfg_orig)


    rollup_config = read_json(paths.rollup_config_path)
    addresses = read_json(paths.addresses_json_path)

    log.info('Bringing up L2.')
    run_command(['docker', 'compose', 'up', '-d', 'l2'], cwd=paths.ops_bedrock_dir, env={
        'PWD': paths.ops_bedrock_dir
    })
    wait_for_rpc_server("http://127.0.0.1:9545")

    log.info('Bringing up everything else.')
    run_command(['docker-compose', 'up', '-d', 'op-node', 'op-proposer', 'op-batcher'], cwd=paths.ops_bedrock_dir, env={
        'PWD': paths.ops_bedrock_dir,
        'L2OO_ADDRESS': addresses['L2OutputOracleProxy'],
        'SEQUENCER_BATCH_INBOX_ADDRESS': rollup_config['batch_inbox_address'],
        'OP_BATCHER_SEQUENCER_BATCH_INBOX_ADDRESS': rollup_config['batch_inbox_address'],
        'INIT_HOLDER_PRV': l1_init_holder_prv,
        'PROPOSER_ADDRESS_PRV': proposer_address_prv
    })

    log.info('Devnet ready.')


def eth_accounts(url):
    log.info(f'Fetch eth_accounts {url}')
    conn = http.client.HTTPConnection(url)
    headers = {'Content-type': 'application/json'}
    body = '{"id":2, "jsonrpc":"2.0", "method": "eth_accounts", "params":[]}'
    conn.request('POST', '/', body, headers)
    response = conn.getresponse()
    data = response.read().decode()
    conn.close()
    return data


def debug_dumpBlock(url):
    log.info(f'Fetch debug_dumpBlock {url}')
    conn = http.client.HTTPConnection(url)
    headers = {'Content-type': 'application/json'}
    body = '{"id":3, "jsonrpc":"2.0", "method": "debug_dumpBlock", "params":["latest"]}'
    conn.request('POST', '/', body, headers)
    response = conn.getresponse()
    data = response.read().decode()
    conn.close()
    return data


def wait_for_rpc_server(url):
    log.info(f'Waiting for RPC server at {url}')
    body = '{"id":1, "jsonrpc":"2.0", "method": "eth_chainId", "params":[]}'
    status = True
    while status:
        try:
            headers = {
                "Content-Type": "application/json"
            }

            response = requests.post(url, headers=headers, data=body)
            if response.status_code != 200:
                time.sleep(5)
            else:
                log.info("Status code is 200, continue next step")
                status = False
        except requests.exceptions.ConnectionError:
                time.sleep(5)


CommandPreset = namedtuple('Command', ['name', 'args', 'cwd', 'timeout'])


def devnet_test(paths):
    # Check the L2 config
    run_command(
        ['go', 'run', 'cmd/check-l2/main.go', '--l2-rpc-url', 'http://localhost:9545', '--l1-rpc-url', 'http://localhost:8545'],
        cwd=paths.ops_chain_ops,
    )

    # Run the two commands with different signers, so the ethereum nonce management does not conflict
    # And do not use devnet system addresses, to avoid breaking fee-estimation or nonce values.
    run_commands([
        CommandPreset('erc20-test',
          ['npx', 'hardhat',  'deposit-erc20', '--network',  'devnetL1',
           '--l1-contracts-json-path', paths.addresses_json_path, '--signer-index', '14'],
          cwd=paths.sdk_dir, timeout=8*60),
        CommandPreset('eth-test',
          ['npx', 'hardhat',  'deposit-eth', '--network',  'devnetL1',
           '--l1-contracts-json-path', paths.addresses_json_path, '--signer-index', '15'],
          cwd=paths.sdk_dir, timeout=8*60)
    ], max_workers=2)


def run_commands(commands: list[CommandPreset], max_workers=2):
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(run_command_preset, cmd) for cmd in commands]

        for future in concurrent.futures.as_completed(futures):
            result = future.result()
            if result:
                print(result.stdout)


def run_command_preset(command: CommandPreset):
    with subprocess.Popen(command.args, cwd=command.cwd,
                          stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True) as proc:
        try:
            # Live output processing
            for line in proc.stdout:
                # Annotate and print the line with timestamp and command name
                timestamp = datetime.datetime.utcnow().strftime('%H:%M:%S.%f')
                # Annotate and print the line with the timestamp
                print(f"[{timestamp}][{command.name}] {line}", end='')

            stdout, stderr = proc.communicate(timeout=command.timeout)

            if proc.returncode != 0:
                raise RuntimeError(f"Command '{' '.join(command.args)}' failed with return code {proc.returncode}: {stderr}")

        except subprocess.TimeoutExpired:
            raise RuntimeError(f"Command '{' '.join(command.args)}' timed out!")

        except Exception as e:
            raise RuntimeError(f"Error executing '{' '.join(command.args)}': {e}")

        finally:
            # Ensure process is terminated
            proc.kill()
    return proc.returncode


def run_command(args, check=True, shell=False, cwd=None, env=None, timeout=None):
    env = env if env else {}
    return subprocess.run(
        args,
        check=check,
        shell=shell,
        env={
            **os.environ,
            **env
        },
        cwd=cwd,
        timeout=timeout
    )


def wait_up(port, retries=10, wait_secs=1):
    for i in range(0, retries):
        log.info(f'Trying 127.0.0.1:{port}')
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        try:
            s.connect(('127.0.0.1', int(port)))
            s.shutdown(2)
            log.info(f'Connected 127.0.0.1:{port}')
            return True
        except Exception:
            time.sleep(wait_secs)

    raise Exception(f'Timed out waiting for port {port}.')

def l1BlockTagGet():
    headers = {
        "Content-Type": "application/json"
    }
    try:
        response = requests.post("http://127.0.0.1:8545",headers=headers,data='{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":74}')
        if response.status_code != 200:
            log.info(f'l1BlockTagGet resp status code is not 200, is {response.status_code}')
            raise Exception("l1BlockTagGet status not 200!")
        else:
            result=response.json()
            log.info(result)
            return result
    except requests.exceptions.ConnectionError:
        raise Exception("l1BlockTagGet connection fail")

def l1BlockTimestampGet(block_tag):
    headers = {
            "Content-Type": "application/json"
    }
    try:
        response = requests.post("http://127.0.0.1:8545",headers=headers,data=f'{{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["{block_tag}", false],"id":74}}')
        if response.status_code != 200:
            log.info(f'l1BlockTimestampGet resp status code is not 200, is {response.status_code}')
            raise Exception("l1BlockTimestampGet status not 200!")
        else:
            return response.json()
    except requests.exceptions.ConnectionError:
        raise Exception("l1BlockTimestampGet connection fail")

def write_json(path, data):
    with open(path, 'w+') as f:
        json.dump(data, f, indent='  ')


def read_json(path):
    with open(path, 'r') as f:
        return json.load(f)
