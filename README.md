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

### Q&A

If you have any questions, you can contact limingxuan2@iie.ac.cn.
