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

    # metadata = {'render.modes': ['human']}   

    def __init__(self, boot=False):
        # 执行一次benchmark，获取所有Prometheus
        self._action_limits = self._get_action_limits()
        # if boot:
        #     self._deploy_cdt()\
        self._init_cdt()
        initial_reward_params = self._collect_state()
        self.initial_reward_params = {
            "Latency": initial_reward_params['caliper']['Latency'],
            "TPS": initial_reward_params['caliper']['TPS']
        }

        self.last_reward_params = {
            "Latency": initial_reward_params['caliper']['Latency'],
            "TPS": initial_reward_params['caliper']['TPS']
        }

        df = pd.DataFrame(initial_reward_params['prom']).astype("float")
        limits = df.mean() * 10
        self.limits = pd.DataFrame(limits).T.to_dict()

        self.obs_num = len(initial_reward_params['prom'][0])
        self.act_num = len(self._action_dict2list(self._action_limits))
        print("obs: %d\t action: %d" % (self.obs_num, self.act_num))
        obs_high = np.ones(self.obs_num, dtype='float32')
        self.observation_space = spaces.Box(
            low= 0,
            high = obs_high,
            dtype = np.float32
        )
        action_high = np.ones(self.act_num, dtype=np.float32)
        self.action_space = spaces.Box(
            low = 0,
            high = action_high,
            dtype = np.float32
        )

        self.c_T = 0.5
        self.c_L = 0.5
        # self.stop_cdt()
        print("initial_reward:\t", self.initial_reward_params)
        print("Aigis init success!")

    def _action_dict2list(self, dict_data):
        res = []
        for item in dict_data:
            for key in dict_data[item]:
                res.append(dict_data[item][key]['value'])
        return res
    def _random_init_action(self, dict_data):
        res = {}
        for item in dict_data:
            for key in dict_data[item]:
                res.append(dict_data[item][key]['value'])
        return res


    def _action_list2dict(self, list_data):
        dict_data = self._action_limits
        index = 0
        for item in dict_data:
            for key in dict_data[item]:
                dict_data[item][key]['value'] = list_data[index]
                index += 1
        return dict_data

    def _get_action_limits(self):
        response = requests.request("GET", url + "/action/limits")
        limits = json.loads(response.text)
        return limits
    
    def _collect_state(self):
        if 'current_reward_params' in dir(self):
            self.last_reward_params = self.current_reward_params

        response = requests.request("GET", url + "/metrics")
        state = json.loads(response.text)
        self.current_reward_params = {
            "Latency": state['caliper']['Latency'],
            "TPS": state['caliper']['TPS']
        }
        self.state = state
        return self.state

    def _cal_reward(self):
        r_T = 0
        r_L = 0
        # deltaT
        deltaT_0 = (self.current_reward_params['TPS'] - self.initial_reward_params['TPS']) / self.initial_reward_params['TPS']
        deltaT_1 = (self.current_reward_params['TPS'] - self.last_reward_params['TPS']) / self.last_reward_params['TPS']
        if deltaT_0 > 0:
            r_T = (np.square((1 + deltaT_0)) -1) * abs(1 + deltaT_1)
        else:
            r_T = -(np.square((1 - deltaT_0)) -1) * abs(1 - deltaT_1)

        # deltaL
        deltaL_0 = (self.current_reward_params['Latency'] - self.initial_reward_params['Latency']) / self.initial_reward_params['Latency']
        deltaL_1 = (self.current_reward_params['Latency'] - self.last_reward_params['Latency']) / self.last_reward_params['Latency']
        if deltaL_0 > 0:
            r_L = (np.square((1 + deltaL_0)) -1) * abs(1 + deltaL_1)
        else:
            r_L = -(np.square((1 - deltaL_0)) -1) * abs(1 - deltaL_1)

        return self.c_L * r_L + self.c_T * r_T

    def _init_cdt(self):
        """
        使用默认参数先跑一次CDT获取状态信息，Reward等
        """
        response = requests.request("POST", url + "/deploy/default")
        print("__init_cdt:\t", response.text)
        time.sleep(1)
        total_iteration = 0
        while not self._check_cdt_status():
            time.sleep(3)
            total_iteration += 1

    def _deploy_cdt(self,config=None):
        """
        config 为ddpg输入，范围0-1
        """
        if config == None:
            # 随机化初始值
            input_config = self._action_dict2list(self._action_limits)
            input_config = list(input_config * np.random.random(len(input_config)))
        else:
            input_list = np.array(config)
            max_list = np.array(self._action_dict2list(self._action_limits))
            input_config = list(input_list * max_list)
            
        print("action: ")
        print(input_config)
        output_config_dict = self._action_list2dict(input_config)
        headers = {
        'Content-Type': 'application/json'
        }
        response = requests.request("POST", url + "/deploy/up", headers=headers, data=json.dumps(output_config_dict))
        # print("deploy_result:\t" + response.text)
        time.sleep(1)
        total_iteration = 0
        while not self._check_cdt_status():
            time.sleep(3)
            total_iteration += 1
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
        return False if response.text == "Not Found" else True

    def _normalization(self, state):
        # min-max 归一化
        df_state = pd.DataFrame(state['prom']).astype('float')
        df_state = pd.DataFrame(df_state.mean()).T
        df_limits = pd.DataFrame(self.limits)

        df_res = df_state / df_limits
        df_res = df_res.fillna(value=0).clip(-1.0, 1.0)
        return np.array(df_res.values.tolist()[0])

    def step(self, action, with_stop = True):
        # 记录起始时间
        start_time = time.time()
        # up
        self._deploy_cdt(list(action))
            
        obs = self._normalization(self._collect_state())
        reward = self._cal_reward()

        if with_stop:
            # down
            self.stop_cdt()

         # 计算时间差值
        seconds, minutes, hours = int(time.time() - start_time), 0, 0

        # 可视化打印
        hours = seconds // 3600
        minutes = (seconds - hours*3600) // 60
        seconds = seconds - hours*3600 - minutes*60
        print("\n  Step complete time cost {:>02d}:{:>02d}:{:>02d}".format(hours, minutes, seconds))
        print("Obs: tps: %s, latency: %s, reward: %f" % (str(self.state['caliper']['TPS']), str(self.state['caliper']['Latency']), reward))

        return obs, reward, True, {}


    def reset(self):
        # 返回obs 状态信息list
        res = self._normalization(self.state)
        # if with_stop:
        #     self.stop_cdt() 
        return res

    def render(self, mode='human'):
        pass
    def close(self):
        self.stop_cdt()