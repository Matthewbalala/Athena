from utils import utils
from collections import defaultdict
import re
import pdb

# peer: { "name": "", "value":"", "unit":"ms"}

class Deployer(object):
    """
    为CDT生成action.yaml
    """
    def __init__(self, default_file="action.default.yaml", 
    source_file="action.max.yaml", target_file="caliper-deploy-tool/action.yaml"):
        self._source_file = source_file
        self._default_file = default_file
        self._target_file = target_file
        self._source_data = utils.load_config(source_file)
        self._internal_max = self._yaml2inter(self._source_data)
        self._default_action = self._yaml2inter(utils.load_config(self._default_file))
    
    def get_limits(self):
        return self._internal_max

    def get_default(self):
        return self._default_action

    def generate(self, config_inter):
        yaml_data = self._inter2yaml(config_inter)
        utils.save_config(yaml_data, self._target_file)


    def _yaml2inter(self, yaml_data):
        res = defaultdict(dict)
        for item in yaml_data:
            for key in yaml_data[item]:
                if type(yaml_data[item][key]) in (int, float):
                    res[item][key] = {
                        "value": yaml_data[item][key],
                        "unit": None
                    }
                else:
                    temp_str = yaml_data[item][key]
                    # print(temp_str)
                    value = re.findall(r'\d', temp_str)
                    index1 = len(temp_str) - len(value)
                    value_num = "".join(value)
                    res[item][key] = {
                        "value": int(value_num),
                        "unit": temp_str[-index1:]
                    }

        # pdb.set_trace()
        return dict(res)


    def _inter2yaml(self, inter_data):
        res = defaultdict(dict)
        CONST_ORDERER_PARAMS = ("BatchTimeout", "MaxMessageCount", "AbsoluteMaxBytes", "PreferredMaxBytes")
        for item in inter_data:
            for key in inter_data[item]:
                temp = int(inter_data[item][key]['value'])
                if key in CONST_ORDERER_PARAMS:
                    temp = temp if temp > 0 else (temp + 1)
                if inter_data[item][key]['unit'] == None:
                    res[item][key] = temp
                else:
                    res[item][key] = str(temp) + inter_data[item][key]['unit']
                # res[item][key] = "{\{action['%s']}\}" % key
        
        # pdb.set_trace()
        return dict(res)

