# 配置
```shell
# 1. main.py
# 修改CDTIP
CDTIP = {your ip}

# 2. caliper-deploy-tool/config.yaml
# 修改config.yaml中的DNS
dnsserver: 10.10.7.51


```

# 运行
```shell
# 1. python main.py

# 2. python train.py
```

# troubleshoot
### parl
```shell
# python3.7
apt install libosmesa6-dev libgl1-mesa-glx libglfw3 # opengl
pip install parl==1.3.1
python -m pip install paddlepaddle==1.8.5 -i https://mirror.baidu.com/pypi/simple
```