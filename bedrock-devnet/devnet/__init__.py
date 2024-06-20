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
import gzip
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
parser.add_argument('--init', help='init BSC L1 devnet chain', type=bool, action=argparse.BooleanOptionalAction)

log = logging.getLogger()

# Global environment variables
DEVNET_NO_BUILD = os.getenv('DEVNET_NO_BUILD') == "true"
DEVNET_FPAC = os.getenv('DEVNET_FPAC') == "true"
DEVNET_PLASMA = os.getenv('DEVNET_PLASMA') == "true"

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
    bsc_dir = pjoin(devnet_dir, 'bsc')
    node_deploy_dir = pjoin(devnet_dir, 'node-deploy')
    node_deploy_genesis_dir = pjoin(node_deploy_dir, 'genesis')
    contracts_bedrock_dir = pjoin(monorepo_dir, 'packages', 'contracts-bedrock')
    deployment_dir = pjoin(contracts_bedrock_dir, 'deployments', 'devnetL1')
    forge_dump_path = pjoin(contracts_bedrock_dir, 'Deploy-900.json')
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
      bsc_dir=bsc_dir,
      node_deploy_dir=node_deploy_dir,
      node_deploy_genesis_dir=node_deploy_genesis_dir,
      contracts_bedrock_dir=contracts_bedrock_dir,
      deployment_dir=deployment_dir,
      forge_dump_path=forge_dump_path,
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

    if args.init:
        bsc_l1_init(paths)
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
    res = eth_accounts('127.0.0.1:8545')

    response = json.loads(res)
    account = response['result'][0]
    log.info(f'Deploying with {account}')

    # send some ether to the create2 deployer account
    run_command([
        'cast', 'send', '--from', account,
        '--rpc-url', 'http://127.0.0.1:8545',
        '--unlocked', '--value', '1ether', '0x3fAB184622Dc19b6109349B94811493BF2a45362'
    ], env={}, cwd=paths.contracts_bedrock_dir)

    # deploy the create2 deployer
    run_command([
      'cast', 'publish', '--rpc-url', 'http://127.0.0.1:8545',
      '0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222'
    ], env={}, cwd=paths.contracts_bedrock_dir)

    fqn = 'scripts/Deploy.s.sol:Deploy'
    run_command([
        'forge', 'script', fqn, '--sender', account,
        '--rpc-url', 'http://127.0.0.1:8545', '--broadcast',
        '--unlocked'
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
    if DEVNET_FPAC:
        deploy_config['useFaultProofs'] = True
        deploy_config['faultGameMaxDuration'] = 10
    if DEVNET_PLASMA:
        deploy_config['usePlasma'] = True
    write_json(paths.devnet_config_path, deploy_config)

def devnet_l1_genesis(paths):
    log.info('Generating L1 genesis state')
    init_devnet_l1_deploy_config(paths)

    fqn = 'scripts/Deploy.s.sol:Deploy'
    run_command([
        'forge', 'script', '--chain-id', '900', fqn, "--sig", "runWithStateDump()"
    ], env={}, cwd=paths.contracts_bedrock_dir)

    forge_dump = read_json(paths.forge_dump_path)
    write_json(paths.allocs_path, { "accounts": forge_dump })
    os.remove(paths.forge_dump_path)

    shutil.copy(paths.l1_deployments_path, paths.addresses_json_path)

def deployL1ContractsForDeploy(paths):
    log.info('Starting L1.')

    run_command(['./start_cluster.sh','start'], cwd=paths.node_deploy_dir)
    wait_up(8545)
    wait_for_rpc_server('http://127.0.0.1:8545')
    time.sleep(3)

    l1env = dotenv_values('./ops-bedrock/l1.env')
    log.info(l1env)
    l1_init_holder = l1env['INIT_HOLDER']
    l1_init_holder_prv = l1env['INIT_HOLDER_PRV']
    proposer_address = l1env['PROPOSER_ADDRESS']
    batcher_address = l1env['BATCHER_ADDRESS']
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

    # send some ether to batcher address
    run_command([
        'cast', 'send', '--private-key', l1_init_holder_prv,
        '--rpc-url', 'http://127.0.0.1:8545', '--gas-price', '10000000000', '--legacy',
        '--value', '10000ether', batcher_address
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

#     log.info('Syncing contracts.')
#     run_command([
#         'forge', 'script', fqn, '--sig', 'sync()',
#         '--rpc-url', 'http://127.0.0.1:8545'
#     ], env={}, cwd=paths.contracts_bedrock_dir)
    log.info('Deployed L1 contracts.')

# Bring up the devnet where the contracts are deployed to L1
def devnet_deploy(paths):
    bsc_l1_init(paths)
    init_devnet_l1_deploy_config(paths)
    l1env = dotenv_values('./ops-bedrock/l1.env')
    log.info(l1env)
    bscChainId = l1env['BSC_CHAIN_ID']
    l1_init_holder = l1env['INIT_HOLDER']
    l1_init_holder_prv = l1env['INIT_HOLDER_PRV']
    proposer_address = l1env['PROPOSER_ADDRESS']
    proposer_address_prv = l1env['PROPOSER_ADDRESS_PRV']
    batcher_address = l1env['BATCHER_ADDRESS']
    batcher_address_prv = l1env['BATCHER_ADDRESS_PRV']
    log.info('Generating network config.')
    devnet_cfg_orig = pjoin(paths.contracts_bedrock_dir, 'deploy-config', 'devnetL1.json')
    devnet_cfg_backup = pjoin(paths.devnet_dir, 'devnetL1.json.bak')
    devnet_cfg_final = pjoin(paths.devnet_dir, 'devnetL1.json')
    shutil.copy(devnet_cfg_orig, devnet_cfg_backup)
    deploy_config = read_json(devnet_cfg_orig)
    deploy_config['l1ChainID'] = int(bscChainId,10)
    deploy_config['l2BlockTime'] = 1
    deploy_config['sequencerWindowSize'] = 14400
    deploy_config['channelTimeout'] = 1200
    deploy_config['l2OutputOracleSubmissionInterval'] = 240
    deploy_config['finalizationPeriodSeconds'] = 3
    deploy_config['enableGovernance'] = False
    deploy_config['eip1559Denominator'] = 8
    deploy_config['eip1559DenominatorCanyon'] = 8
    deploy_config['eip1559Elasticity'] = 2
    deploy_config['l2OutputOracleProposer'] = proposer_address
    deploy_config['batchSenderAddress'] = batcher_address
    deploy_config['baseFeeVaultRecipient'] = l1_init_holder
    deploy_config['l1FeeVaultRecipient'] = l1_init_holder
    deploy_config['sequencerFeeVaultRecipient'] = l1_init_holder
    deploy_config['proxyAdminOwner'] = l1_init_holder
    deploy_config['finalSystemOwner'] = l1_init_holder
    deploy_config['governanceTokenOwner'] = l1_init_holder
    deploy_config['l2GenesisDeltaTimeOffset'] = "0x0"
    deploy_config['fermat'] = 0
    deploy_config['L2GenesisEcotoneTimeOffset'] = "0x0"
    write_json(devnet_cfg_orig, deploy_config)

    if os.path.exists(paths.addresses_json_path):
        log.info('L1 contracts already deployed.')
        log.info('Starting L1.')

        run_command(['./start_cluster.sh','start'], cwd=paths.node_deploy_dir)
        wait_up(8545)
        wait_for_rpc_server('http://127.0.0.1:8545')
    else:
        log.info('Deploying L1 contracts.')
        deployL1ContractsForDeploy(paths)

    l1BlockTag = l1BlockTagGet()["result"]
    log.info(l1BlockTag)
    l1BlockTimestamp = l1BlockTimestampGet(l1BlockTag)["result"]["timestamp"]
    log.info(l1BlockTimestamp)
    deploy_config['l1GenesisBlockTimestamp'] = l1BlockTimestamp
    deploy_config['l1StartingBlockTag'] = l1BlockTag
    write_json(devnet_cfg_orig, deploy_config)
    write_json(devnet_cfg_final, deploy_config)

    if os.path.exists(paths.genesis_l2_path):
        log.info('L2 genesis and rollup configs already generated.')
    else:
        log.info('Generating network config.')
        log.info('Generating L2 genesis and rollup configs.')
        run_command([
            'go', 'run', 'cmd/main.go', 'genesis', 'l2',
            '--l1-rpc', 'http://localhost:8545',
            '--deploy-config', paths.devnet_config_path,
            '--l1-deployments', paths.addresses_json_path,
            '--outfile.l2', paths.genesis_l2_path,
            '--outfile.rollup', paths.rollup_config_path
        ], cwd=paths.op_node_dir)

    if os.path.exists(devnet_cfg_backup):
        shutil.move(devnet_cfg_backup, devnet_cfg_orig)


    rollup_config = read_json(paths.rollup_config_path)
    addresses = read_json(paths.addresses_json_path)

    # Start the L2.
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
        'PROPOSER_ADDRESS_PRV': proposer_address_prv,
        'BATCHER_ADDRESS_PRV': batcher_address_prv
    })

    log.info('Devnet ready.')

def bsc_l1_init(paths):
    l1env = dotenv_values('./ops-bedrock/l1.env')
    log.info(l1env)
    l1_init_holder = l1env['INIT_HOLDER']
    l1_init_holder_prv = l1env['INIT_HOLDER_PRV']
    if os.path.exists(paths.bsc_dir):
        log.info('bsc path exists, skip git clone')
    else:
        run_command(['git','clone','https://github.com/bnb-chain/bsc.git'], cwd=paths.devnet_dir)
        run_command(['git','checkout','v1.4.5'], cwd=paths.bsc_dir)
        run_command(['make','geth'], cwd=paths.bsc_dir)
        run_command(['go','build','-o', './build/bin/bootnode', './cmd/bootnode'], cwd=paths.bsc_dir)
    if os.path.exists(paths.node_deploy_dir):
        log.info('node-deploy path exists, skip git clone')
    else:
        run_command(['git','clone','https://github.com/bnb-chain/node-deploy.git'], cwd=paths.devnet_dir)
        run_command(['git','checkout','27e7ca669a27c8fd259eeb88ba33ef5a1b4ac182'], cwd=paths.node_deploy_dir)
        run_command(['git','submodule','update','--init','--recursive'], cwd=paths.node_deploy_dir)
        run_command(['pip3','install','-r','requirements.txt'], cwd=paths.node_deploy_dir)
        run_command(['npm','install'], cwd=paths.node_deploy_genesis_dir)
        run_command(['forge','install','--no-git','--no-commit','foundry-rs/forge-std@v1.7.3'], cwd=paths.node_deploy_genesis_dir)
        run_command(['poetry','install'], cwd=paths.node_deploy_genesis_dir)
        with open(pjoin(paths.node_deploy_dir,'.env'), 'r') as file:
            file_content = file.read()
        file_content = file_content.replace('0x04d63aBCd2b9b1baa327f2Dda0f873F197ccd186', l1_init_holder)
        file_content = file_content.replace('59ba8068eb256d520179e903f43dacf6d8d57d72bd306e1bd603fdb8c8da10e8', l1_init_holder_prv)
        with open(pjoin(paths.node_deploy_dir,'.env'), 'w') as file:
            file.write(file_content)

    shutil.copy(pjoin(paths.bsc_dir,'build','bin','geth'), pjoin(paths.node_deploy_dir,'bin','geth'))
    shutil.copy(pjoin(paths.bsc_dir,'build','bin','bootnode'), pjoin(paths.node_deploy_dir,'bin','bootnode'))
    if os.path.exists(pjoin(paths.node_deploy_dir,'.local')):
        log.info('already init .local config file, skip init script')
    else:
        run_command(['chmod','+x','start_cluster.sh'], cwd=paths.node_deploy_dir)
        run_command(['./start_cluster.sh','reset'], cwd=paths.node_deploy_dir)
        run_command(['./start_cluster.sh','stop'], cwd=paths.node_deploy_dir)

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
