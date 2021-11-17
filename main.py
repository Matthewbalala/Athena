from flask import Flask, jsonify, request
from utils import utils
from collector import Collector
import pdb
import time
import shutil
import os
from concurrent.futures import ThreadPoolExecutor
import subprocess
import json
from deployer import Deployer

# 全局变量
STATEOFCDT = True
# 创建线程池执行器
executor = ThreadPoolExecutor(2)

action_deployer = Deployer()

CONFIG_PATH = "caliper-deploy-tool/config.yaml"

CDTIP = "10.10.7.51"

app = Flask(__name__)
cdt_config = utils.load_config(CONFIG_PATH)
node_metrics_points = utils.get_node_endpoints(cdt_config)

reportfile = 'caliper-deploy-tool/report.html'
channel_name = 'mychannel'
chaincode_name = 'smallbank'
print("using channe name: {}, chaincode name: {}".format(channel_name, chaincode_name))
metrics_collector = Collector(node_metrics_points, reportfile, channel_name, chaincode_name)



@app.route('/')
def index():
    return "Hello, Aigis!"

# v1 metrics接口
@app.route('/cdt/metrics', methods=['GET'])
def get_metrics():
    # 获取所有peer的平均metrics
    res= {}
    # res['prom'] = metrics_collector.collect_from_prometheus(utils.handler_metrics_prom)
    res['prom'] = metrics_collector.collect_from_prometheus()
    res['caliper'] = metrics_collector.collect_from_caliper()
    # utils.gen_limitscsv(res['prom'])
    return res


@app.route('/cdt/action/status', methods=['GET'])
def get_status():
    global STATEOFCDT
    base_dir = os.path.dirname(__file__)
    target_file = os.path.join(base_dir, "caliper-deploy-tool", "report.html")
    if os.path.exists(target_file):
        return "Exist"
    elif not STATEOFCDT:
        STATEOFCDT = True
        return "Retry"
    else:
        return "Not Found"

@app.route('/cdt/reset', methods=['GET'])
def reset_status():
    global STATEOFCDT
    STATEOFCDT = False
    return "State of cdt: {}\n".format(STATEOFCDT)

@app.route('/cdt/action/limits', methods=['GET'])
def get_limits():
    # pdb.set_trace()
    return action_deployer.get_limits()

@app.route('/cdt/action/ddpg/limits', methods=['GET'])
def get_ddpg_limits():
    # pdb.set_trace()
    return Deployer(source_file="action-ddpg.max").get_limits()

@app.route('/cdt/action/default', methods=['GET'])
def get_default():
    # pdb.set_trace()
    return action_deployer.get_default()
    
# benchmark
def invoke_cdt(auto_stop = False):
    global STATEOFCDT
    STATEOFCDT = True
    if auto_stop:
        shell_cmd_down = "export CDTIP=%s; cd caliper-deploy-tool; make deploy-fabric-down" % CDTIP
        subprocess.call(shell_cmd_down, shell=True)
    else:
        shell_cmd_down = "export CDTIP=%s; cd caliper-deploy-tool; make deploy-fabric-down" % CDTIP
        subprocess.call(shell_cmd_down, shell=True)
        shell_cmd_generate = "export CDTIP=%s; cd caliper-deploy-tool; make generate" % CDTIP
        subprocess.call(shell_cmd_generate, shell=True)
        shell_cmd_up = "export CDTIP=%s; cd caliper-deploy-tool; make deploy-fabric-up" % CDTIP
        subprocess.call(shell_cmd_up, shell=True)
        shell_cmd_start = "export CDTIP=%s; cd caliper-deploy-tool; make start-cdt" % CDTIP
        subprocess.call(shell_cmd_start, shell=True)


@app.route('/cdt/deploy/up', methods=['POST'])
def deploy():
    # 1. 获取配置信息重新生成配置文件
    config = json.loads(request.data.decode())
    print("config data: ", config)
    action_deployer.generate(config)
    # 2. 移动report.html到history文件夹
    base_dir = os.path.dirname(__file__)
    print("Current dir: " + base_dir)
    target_dir = os.path.join(base_dir, "caliper-deploy-tool", "history")
    source_file = os.path.join(base_dir, "caliper-deploy-tool", "report.html")
    print("src: %s,\t target: %s" % (source_file, target_dir))
    if not os.path.exists(target_dir):
        os.mkdir(target_dir)
    try:
        shutil.move(source_file, os.path.join(target_dir, "report-" +str(int(time.time())) + ".html"))
    except IOError as e:
        print(e)
    

    # 3.benchmark
    executor.submit(invoke_cdt)

    return "Good"

@app.route('/cdt/deploy/default', methods=['POST'])
def deploy_default():
    # 1. 获取配置信息重新生成配置文件
    config = utils.load_config("action.default.yaml")
    # pdb.set_trace()
    print(config)
    utils.save_config(config, "caliper-deploy-tool/action.yaml")
    # 2. 移动report.html到history文件夹
    base_dir = os.path.dirname(__file__)
    print("Current dir: " + base_dir)
    target_dir = os.path.join(base_dir, "caliper-deploy-tool", "history")
    source_file = os.path.join(base_dir, "caliper-deploy-tool", "report.html")
    print("src: %s,\t target: %s" % (source_file, target_dir))
    if not os.path.exists(target_dir):
        os.mkdir(target_dir)
    try:
        shutil.move(source_file, os.path.join(target_dir, "report-" +str(int(time.time())) + ".html"))
    except IOError as e:
        print(e)
    

    # 3.benchmark
    executor.submit(invoke_cdt)

    return "Good"

@app.route('/cdt/deploy/down', methods=['DELETE'])
def deploy_down():
    
    # 1. 开始删除容器
    executor.submit(invoke_cdt, True)

    # 2. 移动report.html到history文件夹
    base_dir = os.path.dirname(__file__)
    print("Current dir: " + base_dir)
    target_dir = os.path.join(base_dir, "caliper-deploy-tool", "history")
    source_file = os.path.join(base_dir, "caliper-deploy-tool", "report.html")
    print("src: %s,\t target: %s" % (source_file, target_dir))
    if not os.path.exists(target_dir):
        os.mkdir(target_dir)
    try:
        shutil.move(source_file, os.path.join(target_dir, "report-" +str(int(time.time())) + ".html"))
    except IOError as e:
        print(e)

    return "deploy_down"


if __name__ == '__main__':
    app.run(host = "0.0.0.0", debug=True)