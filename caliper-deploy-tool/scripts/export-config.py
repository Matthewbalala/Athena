import yaml
import os
import pdb
import sys
from jinja2 import Environment,FileSystemLoader
import logging
from collections import defaultdict
logging.basicConfig(format='%(asctime)s - %(pathname)s[line:%(lineno)d] - %(levelname)s: %(message)s',
                    level=logging.DEBUG)

CONFIG_FILE = 'config.yaml'
CONFIG_ACTION = 'action.yaml'
TPL_PEER = 'peer-base.yaml'
TPL_ORDERER = 'orderer-base.yaml'
TPL_CA = 'ca-base.yaml'
TPL_CLIENT = 'client-base.yaml'
DIST_DIR = 'dist'
DNS_CONFIG = 'dns/coredns/hostsfile'
TPL_BENCHMARK = 'benchmarkconfig-base.yaml'
TPL_CONFIGTX = "configtx-base.yaml"
TPL_CRYPTO_CONFIG = "crypto-config-base.yaml"


# load config
def loadconfig(file):
    with open(file) as f:
        data = yaml.load(f, Loader=yaml.FullLoader)
        # print(data)
        # logging.debug(data)
    return data


def merge(config):
    """
    合并yaml,支持单主机部署多个容器
    config -- array

    内容格式:
    [{
        'filename':'fp'
        'yaml':yaml
    },]
    """
    tempdict = defaultdict(list) # key: filename, value: list(yaml)
    def factoryheader():
        return {'version':'3', 'services':{}}
    res = defaultdict(factoryheader)
    addtionaldict = {}
    for item in config:
        if item['filename'] in ('config-distributed.yaml', 'client.yaml', 'configtx.yaml', 'crypto-config.yaml'):
            addtionaldict[item['filename']] = item['yaml']
            continue
        tempdict[item['filename']].append(yaml.load(item['yaml'], Loader=yaml.FullLoader))

    for x in tempdict:
        # peer, orderer, ca
        # logging.debug(x)
        for i in tempdict[x]:
            res[x]['services'] = {**res[x]['services'], **i}
        res[x] = yaml.dump(res[x], default_flow_style=False)

    addtionaldict = {**addtionaldict, **res}

    return addtionaldict

def genCC(config):
    """
    生成configtx.yaml和crypto-config.yaml所需配置
    """
    configtx_config = {}
    crypto_config = {}

    configtx_config["orgs"] = []
    configtx_config["orderers"] = []
    crypto_config["peerorgs"] = []

    orgdict = {}

    crypto_config["orderer_count"] = len(config['fabric-network']['orderer'])

    for org in config['client']['orgs']:
        orgdict["org"] = org
        peer_num = 0
        set_anchor = False
        for peer in config['fabric-network']['peer']:
            if peer.find(org) != -1:
                peer_num += 1
                if not set_anchor:
                    orgdict["anchorpeer"] = {
                        "name": peer,
                        "port": config['fabric-network']['peer'][peer]['port']
                    }
                    set_anchor = True
                    configtx_config["orgs"].append(orgdict)
                    orgdict = {}
        crypto_config["peerorgs"].append({
            "orgname": org + ".example.com",
            "count": peer_num
        })
    for orderer in config['fabric-network']['orderer']:
                configtx_config["orderers"].append({
                    "name": orderer,
                    "port":config['fabric-network']['orderer'][orderer]['port']
                })


    return configtx_config, crypto_config
               



def render(config, dnsfile, action):
    res = []
    tpl_peer = env.get_template(TPL_PEER)
    tpl_orderer = env.get_template(TPL_ORDERER)
    tpl_ca = env.get_template(TPL_CA)
    tpl_client = env.get_template(TPL_CLIENT)
    tpl_benchmark = env.get_template(TPL_BENCHMARK)
    tpl_configtx = env.get_template(TPL_CONFIGTX)
    tpl_crypto = env.get_template(TPL_CRYPTO_CONFIG)
    sources_benchmark = []
    for tp in config['fabric-network']:
        # peer, orderer, ca
        for item in config['fabric-network'][tp]:
            logging.debug(item)
            sources_benchmark.append("http://%s:2375/%s"%(config['fabric-network'][tp][item]['host'], item))
            if item.startswith('peer'):
                res.append({
                    'filename': "docker-compose-%s.yaml"%config['fabric-network'][tp][item]['host'],
                    'yaml': tpl_peer.render(peer = config['fabric-network'][tp][item], name=item, action = action['peer']) 
                })
                sources_benchmark.append("http://%s:2375/dev-%s-smallbank-v0"%(config['fabric-network'][tp][item]['host'], item))
                dnsfile.write(config['fabric-network'][tp][item]['host'] + '   ' + item + '\n')
            if item.startswith('orderer'):
                res.append({
                'filename': "docker-compose-%s.yaml"%config['fabric-network'][tp][item]['host'],
                'yaml': tpl_orderer.render(orderer = config['fabric-network'][tp][item], name=item, action = action['orderer'])
                })
                dnsfile.write(config['fabric-network'][tp][item]['host'] + '   ' + item + '\n')
            if item.startswith('ca'):
                res.append({
                'filename': "docker-compose-%s.yaml"%config['fabric-network'][tp][item]['host'],
                'yaml': tpl_ca.render(ca = config['fabric-network'][tp][item], name=item) 
                })
                dnsfile.write(config['fabric-network'][tp][item]['host'] + '   ' + item + '\n')

    res.append({
        'filename': 'client.yaml',
        'yaml': tpl_client.render(client = config['client'], fabric=config['fabric-network'])
        })

    configtx_conf, crypto_conf = genCC(config)
    
    res.append({
        'filename': 'configtx.yaml',
        'yaml': tpl_configtx.render(orgs=configtx_conf['orgs'], orderers=configtx_conf['orderers'], action = action['configtx'])
        })

    res.append({
        'filename': 'crypto-config.yaml',
        'yaml': tpl_crypto.render(orderer_count=crypto_conf["orderer_count"], peerorgs=crypto_conf["peerorgs"])
        })

    res.append({
        'filename': 'config-distributed.yaml',
        'yaml': tpl_benchmark.render(fabric=sources_benchmark)
        })
    # merge
    res = merge(res)
    dnsfile.close()

    # logging.debug(res)
    # pdb.set_trace()
    return res


def export(arr):
    for item in arr:
        with open(os.path.join(DIST_DIR, item), 'w+') as fout:
            fout.write(arr[item])

def evalgroup(config, groupid="none"):
    """
    按group生成配置
    """
    logging.info("Entered Group Mode!    GroupID: " + groupid)
    drop = []
    for key in config['fabric-network']['peer']:
        if key not in config['caliper-eval']['group'][groupid]:
            drop.append(key)
            
    for key in drop:
        del config['fabric-network']['peer'][key]
    # logging.debug(config['fabric-network']['peer'])
    return config



if __name__ == '__main__':
    # ensure dist
    os.system("mkdir -p %s" % DIST_DIR)
    # python export-config.py group n1
    config = loadconfig(CONFIG_FILE)
    config_action = loadconfig(CONFIG_ACTION)
    if len(sys.argv) == 3:
        # Enter Group Mode 
        config = evalgroup(config, sys.argv[2])

    # pdb.set_trace()
    env = Environment(loader = FileSystemLoader('templates'))
    # add dns 
    dnsfile = open(DNS_CONFIG, 'w')
    arr = render(config, dnsfile, config_action)
    export(arr)

    logging.info("Finished!")

