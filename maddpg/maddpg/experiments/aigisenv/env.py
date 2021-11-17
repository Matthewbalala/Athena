import gym
from gym import error, spaces, utils
from gym.utils import seeding
import numpy as np
import pandas as pd
import requests
import json
import time
import pdb

url = "http://127.0.0.1:5000/cdt"



class AigisEnv(gym.Env):
    """
    Multi Agent
    obs: [[orderer], [peer], [peer-net]]
    action:[orderer, peer, peer-net]
    agent: orderer, peer, peer-net
    """
    # metadata = {'render.modes': ['human']}   

    def __init__(self, booted=True, act_importance=53):
        """
        params:
        booted: 标识是否已经有初始化数据
        act_importance: 调整的action参数数量, 需要加1
        """
        # 重要性排序action列表
        self._internal_actions = [
            'CORE_PEER_GOSSIP_STATE_BLOCKBUFFERSIZE',
            'CORE_PEER_GOSSIP_PUBLISHCERTPERIOD',
            'CORE_PEER_GOSSIP_PROPAGATEITERATIONS',
            'PreferredMaxBytes',
            'CORE_PEER_KEEPALIVE_MININTERVAL',
            'CORE_PEER_DISCOVERY_AUTHCACHEMAXSIZE',
            'CORE_PEER_DELIVERYCLIENT_CONNTIMEOUT',
            'CORE_PEER_KEEPALIVE_CLIENT_INTERVAL',
            'CORE_PEER_CLIENT_CONNTIMEOUT',
            # 'CORE_LEDGER_TOTALQUERYLIMIT',
            'CORE_PEER_GOSSIP_SENDBUFFSIZE',
            'CORE_PEER_KEEPALIVE_CLIENT_TIMEOUT',
            'CORE_PEER_GOSSIP_MAXPROPAGATIONBURSTLATENCY',
            'CORE_PEER_GOSSIP_REQUESTWAITTIME',
            'CORE_PEER_GOSSIP_STATE_MAXRETRIES',
            'CORE_PEER_KEEPALIVE_DELIVERYCLIENT_TIMEOUT',
            'ORDERER_GENERAL_KEEPALIVE_SERVERTIMEOUT',
            'CORE_PEER_GOSSIP_PUBLISHSTATEINFOINTERVAL',
            # 'CORE_CHAINCODE_EXECUTETIMEOUT',
            'CORE_PEER_GOSSIP_MAXBLOCKCOUNTTOSTORE',
            'ORDERER_GENERAL_KEEPALIVE_SERVERMININTERVAL',
            'CORE_PEER_GOSSIP_REQUESTSTATEINFOINTERVAL',
            'CORE_PEER_GOSSIP_ALIVEEXPIRATIONTIMEOUT',
            'CORE_PEER_GOSSIP_STATE_CHECKINTERVAL',
            'CORE_PEER_KEEPALIVE_DELIVERYCLIENT_INTERVAL',
            'CORE_PEER_GOSSIP_RECONNECTINTERVAL',
            'CORE_PEER_GOSSIP_PULLPEERNUM',
            'CORE_PEER_GOSSIP_STATE_BATCHSIZE',
            'CORE_PEER_GOSSIP_RECVBUFFSIZE',
            'CORE_PEER_GOSSIP_RESPONSEWAITTIME',
            'CORE_PEER_GOSSIP_PROPAGATEPEERNUM',
            'CORE_PEER_GOSSIP_PULLINTERVAL',
            'CORE_PEER_GOSSIP_STATE_RESPONSETIMEOUT',
            'AbsoluteMaxBytes',
            'CORE_PEER_GOSSIP_MAXPROPAGATIONBURSTSIZE',
            'BatchTimeout',
            'MaxMessageCount',
            'ORDERER_GENERAL_AUTHENTICATION_TIMEWINDOW',
            'ORDERER_GENERAL_CLUSTER_SENDBUFFERSIZE',
            'ORDERER_GENERAL_KEEPALIVE_SERVERINTERVAL',
            'ORDERER_METRICS_STATSD_WRITEINTERVAL',
            'ORDERER_RAMLEDGER_HISTORYSIZE',
            # 'CORE_METRICS_STATSD_WRITEINTERVAL',
            'CORE_PEER_GOSSIP_MEMBERSHIPTRACKERINTERVAL',
            'CORE_PEER_AUTHENTICATION_TIMEWINDOW',
            'CORE_PEER_DELIVERYCLIENT_RECONNECTTOTALTIMETHRESHOLD',
            'CORE_PEER_DISCOVERY_AUTHCACHEPURGERETENTIONRATIO',
            'CORE_PEER_GOSSIP_ALIVETIMEINTERVAL',
            'CORE_PEER_GOSSIP_CONNTIMEOUT',
            'CORE_PEER_GOSSIP_DIALTIMEOUT',
            'CORE_PEER_GOSSIP_DIGESTWAITTIME',
            'CORE_PEER_GOSSIP_ELECTION_LEADERALIVETHRESHOLD',
            'CORE_PEER_GOSSIP_ELECTION_STARTUPGRACEPERIOD',
            'CORE_PEER_DELIVERYCLIENT_RECONNECTBACKOFFTHRESHOLD',
            'CORE_PEER_GOSSIP_ELECTION_MEMBERSHIPSAMPLEINTERVAL'
        ]

        assert (act_importance - 1) <= len(self._internal_actions)

        # 执行一次benchmark，获取所有Prometheus
        self._act_importance = act_importance
        self._action_limits = self._get_action_limits()
        self._max_list = np.array(self._action_dict2list(self._action_limits, index=True)).copy()
        self._obs_limits = None
        self._obs_shape_dict = None
        

        if not booted:
            self._init_cdt()
            initial_reward_params = self._collect_state()
            
            print("Saved to metadata.")
            self.save_config(initial_reward_params, "/root/code/aigis/maddpg/maddpg/experiments/obs-metadata.npy")
        else:
            print("Loaded from metadata.")
            initial_reward_params = self.load_metadata("/root/code/aigis/maddpg/maddpg/experiments/obs-metadata.npy")
            self.state = initial_reward_params
            self.current_reward_params = {
            "Latency": self.state['caliper']['Latency'],
            "TPS": self.state['caliper']['TPS']
        }

        # set obs limits
        _obs_num, _obs = self._handle_state(initial_reward_params)
        self._init_obs = _obs

        self.initial_reward_params = {
            "Latency": initial_reward_params['caliper']['Latency'],
            "TPS": initial_reward_params['caliper']['TPS']
        }

        self.last_reward_params = {
            "Latency": initial_reward_params['caliper']['Latency'],
            "TPS": initial_reward_params['caliper']['TPS']
        }


        # df = pd.DataFrame(initial_reward_params['prom']).astype("float")
        # limits = df.mean() * 10
        # self.limits = pd.DataFrame(limits).T.to_dict()

        # self.obs_num = len(initial_reward_params['prom'][0])
        # self.obs_num = 54
        self.n = len(_obs)
        print("number of agent: ", self.n)
        # self.act_num = len(self._action_dict2list(self._action_limits))
        print("obs: %d\t current action: %d\t total action: %d" % (_obs_num, 
        len(self._act_range_dict['orderer']) + len(self._act_range_dict['peer']) + len(self._act_range_dict['peer-net']),
        self._act_range_dict['length']))
        # obs_high = [np.ones(i.shape, dtype='float32') for i in _obs]
        # print(obs_high)
        # self.observation_space = [spaces.Box(
        #     low= 0.0,
        #     high = 1.0,
        #     shape = i.shape,
        #     dtype = np.float32
        # ) for i in _obs]
        # action_high = np.ones(self.act_num, dtype=np.float32)
        self.action_space = [
            spaces.Box(
            low = 0.0,
            high = 1.0,
            shape = (len(self._act_range_dict['orderer']),),
            dtype = np.float32
            ),
            spaces.Box(
            low = 0.0,
            high = 1.0,
            shape = (len(self._act_range_dict['peer']),),
            dtype = np.float32
            ),
            spaces.Box(
            low = 0.0,
            high = 1.0,
            shape = (len(self._act_range_dict['peer-net']),),
            dtype = np.float32
            )
        ]
        self.observation_space = [i.shape for i in _obs]
        # self.action_space = self.act_num


        self.err_count = 0
        self.c_T = 0.8
        self.c_L = 0.2

        

        # self.stop_cdt()
        print("initial_reward:\t", self.initial_reward_params)
        print("Aigis init success!")

    def _convert2sorteddf(self, data):
        """
        处理原始state信息，返回《固定排序》的DataFrame
        """
        tdf = pd.DataFrame.from_dict(data).fillna(0).astype("float")
        tcols = sorted(list(tdf.columns))
        return tdf[tcols]
    
    def _handle_state(self, state):
        """
        26, 25, 29
        处理原始state信息，转换成二维obs[[orderer], [peer], [peer-net]]
        """
        _obs = []
        for node_name in ("orderer", "peer", "peer-net"):
            _obs.append(pd.DataFrame(self._convert2sorteddf(state['prom'][node_name]).mean()).T)
        _onelinedf = pd.concat(_obs, axis=1)
        if self._obs_limits is  None:
            # 设置obs的limit为初始值的5倍
            self._obs_limits = _onelinedf * 5
        
        # 此处fillna是为了防止除数为0，即最大值有可能是0
        _onelineres = (_onelinedf / self._obs_limits).fillna(1.0).clip(0, 1.0)
        # print(_onelineres.info())
        _res = []
        _sum = 0
        for item in _obs: 
            _res.append(np.squeeze(_onelineres.iloc[:, _sum: item.shape[1] + _sum].values, 0))
            _sum += item.shape[1]
        
        if self._obs_shape_dict is None:
            self._obs_shape_dict = {}
            self._obs_shape_dict["orderer"] = _res[0].shape
            self._obs_shape_dict["peer"] = _res[1].shape
            self._obs_shape_dict["peer-net"] = _res[2].shape
            print("set obs_shape_dict: {}, {}, {}".format(_res[0].shape, _res[1].shape, _res[2].shape))

        return _onelineres.shape[1], _res



    # load config
    def load_metadata(self, file):
        return np.load(file, allow_pickle=True).item()

    # save config
    def save_config(self, data, target_file):
        np.save(target_file, data)

    def _action_dict2list(self, dict_data, index=False):
        """
        index: 是否需要计算索引
        input: ["configtx", "orderer", "peer"]
                    convert
        return: list [[orderer], [peer], [peer-net]]
        """
        res = []
        # 记录输出的action在原始一维action数值中的位置【对应action_limits】
        _index_range = {
            "orderer": [],
            "peer" : [],
            "peer-net": [],
            "length": 0 # length为总action的长度，但可能不等于len(orderer + peer + peer-net)
        }
        # for item in dict_data:
        #     node_acts = []
        #     for key in dict_data[item]:
        #         node_acts.append(dict_data[item][key]['value'])
        #     res.append(node_acts)
        
        # 标识gossip的action在action一维数组的索引
        flag = 0
        # orderer
        orderer_acts = []
        for item in ('configtx', 'orderer'):
            for key in dict_data[item]:
                if key in self._internal_actions[:self._act_importance]:
                    _index_range['orderer'].append(flag)
                    orderer_acts.append(dict_data[item][key]['value'])
                flag += 1
        res.append(orderer_acts)

        # peer
        peer_acts = []
        net_acts = []

        for key in dict_data['peer']:
            if key in self._internal_actions[:self._act_importance]:
                if "GOSSIP" in key:
                    
                    _index_range['peer-net'].append(flag)
                    net_acts.append(dict_data['peer'][key]['value'])
                else:
                    
                    _index_range['peer'].append(flag)
                    peer_acts.append(dict_data['peer'][key]['value'])
            flag += 1
        
        # 需要保证action的数量使得agent==3
        assert len(peer_acts) != 0
        assert len(net_acts) != 0
        assert len(orderer_acts) != 0

        res.append(peer_acts)
        res.append(net_acts)
        _index_range['length'] = flag

        if index:
            self._act_range_dict = _index_range
            print("set internal action range dict: ", self._act_range_dict)
        return res



    def _action_list2dict(self, list_data):
        """
        input: list [[orderer], [peer], [peer-net]]
                    convert
        return: ["configtx", "orderer", "peer"]
        """
        # 二维变为一维
        flatten_list = [0] * self._act_range_dict['length']
        # 提前赋值
        dict_data = self._action_limits.copy()
        index = 0
        for item in dict_data:
            for key in dict_data[item]:
                flatten_list[index] = dict_data[item][key]['value']
                index += 1
        # orderer
        for (index, item) in enumerate(list_data[0]):
            flatten_list[self._act_range_dict['orderer'][index]] = item
        # peer
        for (index, item) in enumerate(list_data[1]):
            flatten_list[self._act_range_dict['peer'][index]] = item
        # net
        for (index, item) in enumerate(list_data[2]):
            flatten_list[self._act_range_dict['peer-net'][index]] = item


        # dict_data = self._action_limits.copy()
        # index = 0
        # for item in dict_data:
        #     for key in dict_data[item]:
        #         dict_data[item][key]['value'] = flatten_list[index]
        #         index += 1
        return dict_data

    def _get_action_limits(self):
        response = requests.request("GET", url + "/action/limits")
        limits = json.loads(response.text)
        return limits
    
    def _collect_state(self):
        if 'current_reward_params' in dir(self):
            print("Set last_reward_params:\t", self.current_reward_params)
            self.last_reward_params = self.current_reward_params

        response = requests.request("GET", url + "/metrics")
        state = json.loads(response.text)
        self.current_reward_params = {
            "Latency": state['caliper']['Latency'],
            "TPS": state['caliper']['TPS']
        }
        print("Set current_reward_params:\t", self.current_reward_params)
        self.state = state
        return self.state


    # def _cal_delta(self, delta0, delta1):

    #     if delta0 > 0:
    #         _res = (np.square((1 + delta0)) -1) * abs(1 + delta1)
    #     else:
    #         _res = -(np.square((1 - delta0)) -1) * abs(1 - delta1)

    #     if _res > 0 and delta1 < 0:
    #         _res = 0
    #     return _res
    def _cal_delta_T(self, delta0, delta1, eta):
        if delta1 > 0:
            return np.exp(eta * delta0*delta1)
        else:
            return -np.exp(-eta*delta0*delta1)
    
    def _cal_delta_L(self, delta0, delta1, eta):
        if delta1 > 0:
            return -np.exp(eta*delta0*delta1)
        else:
            return np.exp(-eta*delta0*delta1)


    def _cal_reward(self):
        r_T = 0
        r_L = 0
        # deltaT
        deltaT_0 = (self.current_reward_params['TPS'] - self.initial_reward_params['TPS']) / self.initial_reward_params['TPS']
        deltaT_1 = (self.current_reward_params['TPS'] - self.last_reward_params['TPS']) / self.last_reward_params['TPS']
        # if deltaT_0 > 0:
        #     r_T = (np.square((1 + deltaT_0)) -1) * abs(1 + deltaT_1)
        # else:
        #     r_T = -(np.square((1 - deltaT_0)) -1) * abs(1 - deltaT_1)
        r_T = self._cal_delta_T(deltaT_0, deltaT_1, 10)
        # r_T = self._cal_delta(deltaT_0, deltaT_1)

        # deltaL
        deltaL_0 = -(self.current_reward_params['Latency'] - self.initial_reward_params['Latency']) / self.initial_reward_params['Latency']
        deltaL_1 = -(self.current_reward_params['Latency'] - self.last_reward_params['Latency']) / self.last_reward_params['Latency']
        # if deltaL_0 > 0:
        #     r_L = (np.square((1 + deltaL_0)) -1) * abs(1 + deltaL_1)
        # else:
        #     r_L = -(np.square((1 - deltaL_0)) -1) * abs(1 - deltaL_1)
        r_L = self._cal_delta_L(deltaL_0, deltaL_1, 10)
        # r_L = self._cal_delta(deltaL_0, deltaL_1)
        res = self.c_L * r_L + self.c_T * r_T
        return res


    def _init_cdt(self):
        """
        使用默认参数先跑一次CDT获取状态信息，Reward等
        """
        response = requests.request("POST", url + "/deploy/default")
        print("__init_cdt:\t", response.text)
        time.sleep(5)
        while True:
            scode = self._check_cdt_status()
            if scode == 1:
                # 重试 Retry
                return False
            elif scode == 0:
                # 成功
                break
            time.sleep(3)
        print("Set default state successfully.")
        return True


    def _deploy_cdt(self,config=None):
        """
        config 为ddpg输入，范围0-1
        """
        print("_deploy_cdt: config", config)
        time_start = time.time()

        # input_list = np.array(config)
        # input_config = list(input_list * self._max_list)
        input_config = []
        for (index, item) in enumerate(config):
            input_config.append(self._max_list[index] * item)
            
        # print("action: ")
        # print(input_config)
        # return
        output_config_dict = self._action_list2dict(input_config)
        headers = {
        'Content-Type': 'application/json'
        }
        time_end = time.time()
        print("generate action time cost: {}s".format(time_end - time_start))
        response = requests.request("POST", url + "/deploy/up", headers=headers, data=json.dumps(output_config_dict))
        # print("deploy_result:\t" + response.text)
        time.sleep(3)
        # total_iteration = 0
        while True:
            scode = self._check_cdt_status()
            if scode == 1:
                # 重试
                return False
            elif scode == 0:
                # 成功
                break
            time.sleep(3)
            # total_iteration += 1
        
        return True
            # if(total_iteration % 10 == 0):
                # print("[deploy-cdt-up] Wait for cdt benchmark result ...")
        # return True if response.text == "Good" else False

    def stop_cdt(self):
        response = requests.request("DELETE", url + "/deploy/down")
        # print(response.text)
        time.sleep(1)
        total_iteration = 0
        while self._check_cdt_status():
            time.sleep(1)
            total_iteration += 1
            # if(total_iteration % 10 == 0):
                # print("[deploy-cdt-down] Wait for cdt stop ...")
        # print("Done.")
        return True

    
    def _check_cdt_status(self):
        response = requests.request("GET", url + "/action/status")
        if response.text == "Not Found":
            return -1
        elif response.text == "Retry":
            return 1
        else:
            return 0

    # def _normalization(self, state):
    #     # min-max 归一化
    #     df_state = pd.DataFrame(state['prom']).astype('float')
    #     df_state = pd.DataFrame(df_state.mean()).T
    #     df_limits = pd.DataFrame(self.limits)

    #     df_res = df_state / df_limits
    #     df_res = df_res.fillna(value=0).clip(0, 1.0)
    #     _data = np.array(df_res.values.tolist()[0])
    #     print("prom: \t{}".format(_data))
    #     # df_res = df_state.fillna(value=0)
    #     # _data = np.array(df_res.values.tolist()[0])
    #     # _data = (_data - np.mean(_data)) / np.std(_data)
    #     return _data

    def step(self, action, with_stop = True):
        # 记录起始时间
        start_time = time.time()
        # up
        nisRetry = self._deploy_cdt(list(action))

        if nisRetry:
            _, obs = self._handle_state(self._collect_state())
            reward = self._cal_reward()
        else:
            print("ERROR: CDT is down.")
            # 重试时说明action错误导致fabric网络没有正常启动
            self.state['caliper']['TPS'] = 1.0
            self.state['caliper']['Latency'] = 100
            self.current_reward_params = {
            "Latency": self.state['caliper']['Latency'],
            "TPS": self.state['caliper']['TPS']
            }
            if 'current_reward_params' in dir(self):
                print("Set last_reward_params:\t", self.current_reward_params)
                self.last_reward_params = self.current_reward_params
            # _, obs = self._handle_state(self.state)
            # reward = self._cal_reward() - self.err_count * 1.0
            obs = []
            for n in ("orderer", "peer", "peer-net"):
                obs.append(np.zeros(self._obs_shape_dict[n]))
            reward = -100.0
            # self.err_count += 1.0

        # if with_stop:
        #     # down
        #     self.stop_cdt()

         # 计算时间差值
        seconds, minutes, hours = int(time.time() - start_time), 0, 0

        # 可视化打印
        hours = seconds // 3600
        minutes = (seconds - hours*3600) // 60
        seconds = seconds - hours*3600 - minutes*60
        print("[Env] step complete time cost {:>02d}:{:>02d}:{:>02d}".format(hours, minutes, seconds))
        print("[Env] obs:{} tps: {}, latency: {}, reward: {}".format( obs, str(self.state['caliper']['TPS']), str(self.state['caliper']['Latency']), reward))
        for (index, item) in  enumerate(obs):
            print("[ENV] obs shape: [{}]  {}".format(index, item.shape))
        
        assert obs[0].shape == self._obs_shape_dict['orderer']
        assert obs[1].shape == self._obs_shape_dict['peer']
        assert obs[2].shape == self._obs_shape_dict['peer-net']
        
        # print("[Env] obs: prom: {}".format(self.state['prom']))
        return obs, [reward] * self.n, False, {}, self.state['caliper']['TPS'], self.state['caliper']['Latency']


    def reset(self):
        # 返回obs 状态信息list
        return self._init_obs

    def render(self, mode='human'):
        pass
    def close(self):
        self.stop_cdt()


if __name__ == '__main__':
    # aigis = AigisEnv(booted=True, act_importance=53)
    aigis = AigisEnv(booted=False)
    print(aigis.observation_space)
    print(aigis.n)
    print(aigis.action_space)
