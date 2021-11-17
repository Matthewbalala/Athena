# Caliper-Deploy-Tool(CDT)

> 一个基于Hyperledger Caliper v0.3.2的fabric网络部署和测试工具，具有以下特性：
>
> - 部署简化，docker一键安装caliper
> - 支持使用ansible多节点部署fabric 网络，包括docker镜像分发，docker daemon配置修改，启动和停止fabric网络，nfs共享，冗余文件清理
> - 配置文件模板化，简化多节点配置文件生成过程
>
> - 当前支持fabric v1.4.4 raft共识算法

### Runtime

- OS - Ubuntu 16.04 ^   
- OS User: root
- python3 & pip3
- docker-compose
- docker API version: 1.12^
- chrony [optional, 用来同步时间,其他NTP也可]

### Install

- ansible

  ```shell
  pip3 install ansible
  ```

- jinja2

  ```shell
  pip3 install jinja2
  ```

  安装完成。接下来可以准备运行了。

# Quick Start

> 使用默认设置快速开始。多机运行模式（所有主机符合Runtime要求），使用单台主机运行CDT docker容器，其他远程主机运行fabric网络【raft， 3 orderer nodes, 2 org, 4 peer nodes, 2 ca nodes， 智能合约使用smallbank】。

- 环境配置

  | 主机名 | IP           | 角色                                          |
  | ------ | ------------ | --------------------------------------------- |
  | 主机A  | 10.10.28.116 | CDT                                           |
  | 主机B  | 10.10.7.46   | peer0.org1, peer1.org1                        |
  | 主机C  | 10.10.7.47   | peer0.org2, peer1.org2                        |
  | 主机D  | 10.10.7.52   | ca.org1, ca.org2, orderer0, ordere1, orderer2 |

  1. ssh免密登录

  将主机A的ssh公钥添加到B、C、D主机上，并配置B、C、D主机的ssh免密登录。

  

  2. ansible主机设置

     ```shell
     vim /etc/ansible/hosts
     
     # 输入以下内容
     [cdt]
     10.10.7.46 ip=10.10.7.46
     10.10.7.47 ip=10.10.7.47
     10.10.7.52 ip=10.10.7.52
     
     # 测试连通性
     ansible cdt -m ping
     ```

     

  

  3. 配置主机A、B、C、D的环境

     以下所有操作都在主机A上执行, 工作目录为caliper-deploy-tool, 需要关闭systemd-redsolved服务：

      ```shell
      # ensure 127.0.0.1 as host A's dns
      # In ubuntu, you can edit `/etc/resolv.conf` to set "nameserver" to "127.0.0.1.53"

      # setup enviroments: 
      # 1. build cdt docker
      # 2. stop systemd-resolved and release dns port 53
      # ensure 127.0.0.1 as host A's dns
      # In ubuntu, you can edit `/etc/resolv.conf` to set "nameserver" to "127.0.0.1.53"
      # 3. boot coredns server, 
      # 4. boot nfs server using docker 
      make setup-cdt

      # download fabric docker image v1.4.4 and distribute them on host B、C、D [镜像分发]
      make distribute-docker-images
     
      # setup fabric nodes host B、C、D
      # !!! 这个操作会关闭目标主机B、C、D的docker服务，并开放docker服务的远程访问权限、替换daemon.json和docker service
      # 1. update docker config [use reomte docker stat api]
      # 2. apt install nfs-common
      # 3. mount 10.10.28.116:/root/ansible /root/ansible
      export CDTIP=10.10.28.116  # set cdtip accodring to your CDT host
      make setup-config
      ```

- 启动CDT benchmark进行测试

  - 修改CDT配置config.yaml 

    有关config.yaml文件的说明请查看#【config.yaml说明】

    ```yaml
    common: &common
      version: cdt
      network: ansible_default
      mountpath: /root/ansible/nfs
      dnsserver: 10.10.28.116
    
    client:
      cryptopath: benchmarks/config_raft
      usecmd: false
      orgs: 
        - org1
        - org2
    fabric-network:
      peer:
        peer0.org1.example.com:  
          <<: *common
          host: 10.10.7.46
          port: 7051
        peer1.org1.example.com:  
          <<: *common
          host: 10.10.7.46
          port: 7061
        peer0.org2.example.com:  
          <<: *common
          host: 10.10.7.47
          port: 7051
        peer1.org2.example.com:  
          <<: *common
          host: 10.10.7.47
          port: 7061
    
      orderer:
        orderer0.example.com:  
          <<: *common
          host: 10.10.7.52
          port: 7050
        orderer1.example.com:  
          <<: *common
          host: 10.10.7.52
          port: 8050
        orderer2.example.com:  
          <<: *common
          host: 10.10.7.52
          port: 9050
      ca:
        ca.org1.example.com:  
          <<: *common
          host: 10.10.7.52
          port: 7054
        ca.org2.example.com:  
          <<: *common
          host: 10.10.7.52
          port: 8054
    ```

    

  - 生成caliper和fabric配置

    ```shell
    # operations:
    # 1. generate dist config
    # 2. generate fabric crypto-config
    make generate
    ```

  - 启动fabric网络

    ```shell
    # use ansible to boot fabric network on host B、C、D
    make deploy-fabric-up
    ```

    

  - 启动CDT进行一次benchmark

    ```shell
    export CDTIP=10.10.28.116
    make start-cdt
    ```

- 测试完成后【环境清理】

  1. 关闭fabric网络

     ```shell
     # use ansible to close fabric docker containers on host B、C、D
     make deploy-fabric-down
     ```

     

  2. [可选]重置环境

  ```shell
  # 1. rm cdt docker image.
  # 2. stop&rm coredns
  # 3. umount nfs
  # 4. stop&rm nfs
  # 5. rm -rf /root/ansible  on each host
  # 6. stop&rm fabric's docker container on each host
  # 7. (non-implements) restore docker config on each host
  make purge
  ```

  

# 运行结果

CDT运行结果存放在result文件夹，文件格式为html。

# config.yaml说明

```yaml
common: &common											# fabric节点的公共配置
  version: cdt											# CDT分发的fabric镜像tag名称，不支持修改
  network: ansible_default								# docker-compose文件使用的网络
  mountpath: /root/ansible/nfs							# fabric网络运行节点的挂载点，用于从CDT节点共享fabric MSP 证书
  dnsserver: 10.10.28.116								# CDT节点IP，运行coredns服务器
	
client:													# caliper 客户端配置
  cryptopath: benchmarks/config_raft					# 客户端证书路径前缀
  usecmd: false											# 客户端是否使用启动和结束命令，不支持修改
  orgs:													# 客户端组织数组
    - org1
    - org2
    
fabric-network:											# fabric网络配置
  peer:													# peer节点
    peer0.org1.example.com:  
      <<: *common
      host: 10.10.7.46									# 主机IP
      port: 7051										# 对外暴露端口，此端口映射与宿主机上。【同一主机端口设置不能相同】
    peer1.org1.example.com:  
      <<: *common
      host: 10.10.7.46
      port: 7061
  
  orderer:												# orderer节点
    orderer0.example.com:  
      <<: *common
      host: 10.10.7.52
      port: 7050

  ca:													# ca节点
    ca.org1.example.com:  
      <<: *common
      host: 10.10.7.52
      port: 7054
```

