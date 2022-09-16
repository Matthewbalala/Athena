# Athena
> This tool is an implementation in the paper "Auto-Tuning with Reinforcement Learning for Permissioned-Blockchain Systems".

### Environment
- caliper-deploy-tool
```shell
# https://github.com/konoleoda/caliper-deploy-tool
# Install according to CDT's Readme and perform Quick Start
```
- Athena
```shell
pip install -r requirements.txt
```
### RUN
Start two terminals and execute the following commands respectively.
```shell
# 1. Opening the parameter adjustment servers
python main.py
# 2. Training the model.
cd maddpg/maddpg/experiments && python train.py
```

### Note
Parameters list and parameters rank in branch Athena-Fabric-v2.4.3.

We have supported Fabric v1.4.x and Fabric v2.x.x , you can use it by switching branch. 

v2.4.3 of FastFabric is in branch FastFabric.

### Q&A

If you have any questions, you can contact limingxuan2@iie.ac.cn.
