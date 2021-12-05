# Athena
> This tool is an implementation in the paper "Auto-Tuning with Reinforcement Learning for Permissioned-Blockchain Systems".

### 环境配置
- caliper-deploy-tool
```shell
# https://github.com/konoleoda/caliper-deploy-tool
# 根据CDT的Readme进行安装，并进行Quick Start
```
- Athena
```shell
pip install -r requirements.txt
```
### 运行
分别启动两个终端执行
```shell
# 1. 
python main.py
# 2. 
cd maddpg/maddpg/experiments && python train.py
```

