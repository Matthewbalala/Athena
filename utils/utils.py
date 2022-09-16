import yaml
import pandas as pd
import pdb



# load config
def load_config(file):
    with open(file) as f:
        data = yaml.load(f, Loader=yaml.FullLoader)
        # print(data)
        # logging.debug(data)
    return data

# save config
def save_config(yaml_data, target_file):
    with open(target_file, "w", encoding="utf-8") as f:
            yaml.dump(yaml_data, f)

def get_peer_mpoint(config_data):
    res = []
    for i in config_data['fabric-network']['peer']:
        res.append("http://" + config_data['fabric-network']['peer'][i]['host']
        + ":" + str(config_data['fabric-network']['peer'][i]['port'] + 1) + "/" + "metrics")

    return res

def get_node_endpoints(config_data):
    res = [
    ]
    for node_type in ("peer", "orderer"):
        for i in config_data['fabric-network'][node_type]:
            res.append({"url" : "http://" + config_data['fabric-network'][node_type][i]['host']
            + ":" + str(config_data['fabric-network'][node_type][i]['port'] + 1) + "/" + "metrics", "node_type": node_type})

    return res

def handler_metrics_prom(data):
    # 取平均值
    df = pd.DataFrame(data)
    return df.mean().tolist()

# def gen_limitscsv(data):
#     df = pd.DataFrame(data).astype("float")
#     limits = df.mean() * 10
#     limits = pd.DataFrame(limits).T
#     limits.to_csv("gym-aigis/gym_aigis/envs/limits.csv", index=False) 

# def cal_limits_metrics(data):
#     df = pd.DataFrame(data).astype("float")
#     limits = df.mean() * 10
#     limits = pd.DataFrame(limits).T
#     return limits.to_dict()