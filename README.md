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
# 1. 
python main.py
# 2. 
cd maddpg/maddpg/experiments && python train.py
```

### Q&A

If you have any questions, you can contact limingxuan2@iie.ac.cn

### 参数列表

| **Parameter Name**                                             | **Description**                                                                                                                | **Range**        | **Type** |
|----------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------|------------------|----------|
| AbsoluteMaxBytes                                               | The absolute maximum number of bytes allowed for the serialized messages in a batch\.                                          | 512 kb\-10240 kb | orderer  |
| BatchTimeout                                                   | The amount of time to wait before creating a batch\.                                                                           | 100\-2000 ms     | orderer  |
| MaxMessageCount                                                | The maximum number of messages to permit in a batch\.                                                                          | 250\-5000        | orderer  |
| PreferredMaxBytes                                              | The preferred maximum number of bytes allowed for the serialized messages in a batch\.                                         | 1024 kb\-2048 kb | orderer  |
| ORDERER\_GENERAL\_AUTHENTICATION\_TIMEWINDOW             | The acceptable difference between the current server time and the client's time as specified in a client request message\.     | 0\.75\-15 m      | orderer  |
| ORDERER\_GENERAL\_CLUSTER\_SENDBUFFERSIZE                | The maximum number of messages in the egress buffer\.                                                                          | 0\.5\-10         | orderer  |
| ORDERER\_GENERAL\_KEEPALIVE\_SERVERINTERVAL              | The minimum permitted time between client pings\.                                                                              | 360\-7200 s      | orderer  |
| ORDERER\_GENERAL\_KEEPALIVE\_SERVERMININTERVAL           | The time between pings to clients\.                                                                                            | 3\-60 s          | orderer  |
| ORDERER\_GENERAL\_KEEPALIVE\_SERVERTIMEOUT               | The duration the server waits for a response from a client before closing the connection\.                                     | 1\-20 s          | orderer  |
| ORDERER\_METRICS\_STATSD\_WRITEINTERVAL                  | The interval at which locally cached counters and gauges are pushed to statsd; timings are pushed immediately\.                | 1\.5\-30 s       | orderer  |
| ORDERER\_RAMLEDGER\_HISTORYSIZE                            | The number of blocks that the RAM ledger is set to retain\.                                                                    | 50\-1000         | orderer  |
| CORE\_PEER\_KEEPALIVE\_MININTERVAL                       | The minimum permitted time between client pings\.                                                                              | 3\-60 s          | peer     |
| CORE\_PEER\_KEEPALIVE\_CLIENT\_INTERVAL                | The time between pings to peer nodes\.                                                                                         | 3\-60 s          | peer     |
| CORE\_PEER\_KEEPALIVE\_CLIENT\_TIMEOUT                 | The duration the client waits for a response from  peer nodes before closing the connection\.                                  | 1\-20 s          | peer     |
| CORE\_PEER\_KEEPALIVE\_DELIVERYCLIENT\_INTERVAL        | The time between pings to ordering nodes\.                                                                                     | 3\-60 s          | peer     |
| CORE\_PEER\_KEEPALIVE\_DELIVERYCLIENT\_TIMEOUT         | The duration the client waits for a response from ordering nodes before closing the connection\.                               | 1\-20 s          | peer     |
| CORE\_PEER\_AUTHENTICATION\_TIMEWINDOW                   | The acceptable difference between the current server time and the client's time as specified in a client request message       | 0\.75\-15 ms     | peer     |
| CORE\_PEER\_CLIENT\_CONNTIMEOUT                          | Connection timeout\.                                                                                                           | 0\.15\-3 s       | peer     |
| CORE\_PEER\_DELIVERYCLIENT\_RECONNECTTOTALTIMETHRESHOLD  | The total time to spend retrying connections to ordering nodes before giving up and returning an error\.                       | 180\-3600 s      | peer     |
| CORE\_PEER\_DELIVERYCLIENT\_CONNTIMEOUT                  | The connection timeout when connecting to ordering service nodes\.                                                             | 0\.15\-3 s       | peer     |
| CORE\_PEER\_DELIVERYCLIENT\_RECONNECTBACKOFFTHRESHOLD    | The maximum delay between consecutive connection retry attempts ordering nodes\.                                               | 180\-3600 s      | peer     |
| CORE\_PEER\_DISCOVERY\_AUTHCACHEMAXSIZE                  | The maximum size of the cache, after which a purge takes place\.                                                               | 50\-1000         | peer     |
| CORE\_PEER\_DISCOVERY\_AUTHCACHEPURGERETENTIONRATIO      | The proportion \(0 to 1\) of entries that remain in the cache after the cache is purged due to overpopulation\.                | 0\.0375\-0\.75   | peer     |
| CORE\_PEER\_GOSSIP\_MEMBERSHIPTRACKERINTERVAL            | Interval for membershipTracker polling\.                                                                                       | 0\.25\-5 s       | network  |
| CORE\_PEER\_GOSSIP\_MAXBLOCKCOUNTTOSTORE                 | Maximum count of blocks stored in memory\.                                                                                     | 5\-100           | network  |
| CORE\_PEER\_GOSSIP\_MAXPROPAGATIONBURSTLATENCY           | Max time between consecutive message pushes\.                                                                                  | 0\.2\-10 ms      | network  |
| CORE\_PEER\_GOSSIP\_MAXPROPAGATIONBURSTSIZE              | Max number of messages stored until a push is triggered to remote peers\.                                                      | 0\.2\-10         | network  |
| CORE\_PEER\_GOSSIP\_PROPAGATEITERATIONS                  | Number of times a message is pushed to remote peers\.                                                                          | 0\.02\-1         | network  |
| CORE\_PEER\_GOSSIP\_PROPAGATEPEERNUM                     | Number of peers selected to push messages to                                                                                   | 0\.15\-3         | network  |
| CORE\_PEER\_GOSSIP\_PULLINTERVAL                         | Determines frequency of pull phases\.                                                                                          | 0\.2\-4 s        | network  |
| CORE\_PEER\_GOSSIP\_PULLPEERNUM                          | Number of peers to pull from\.                                                                                                 | 0\.15\-3         | network  |
| CORE\_PEER\_GOSSIP\_REQUESTSTATEINFOINTERVAL             | Determines frequency of pulling state info messages from peers\.                                                               | 0\.2\-4 s        | network  |
| CORE\_PEER\_GOSSIP\_PUBLISHSTATEINFOINTERVAL             | Determines frequency of pushing state info messages to peers\.                                                                 | 0\.2\-4 s        | network  |
| CORE\_PEER\_GOSSIP\_PUBLISHCERTPERIOD                    | Time from startup certificates are included in Alive messages\.                                                                | 0\.5\-10 s       | network  |
| CORE\_PEER\_GOSSIP\_DIALTIMEOUT                          | Dial timeout\.                                                                                                                 | 0\.15\-3 s       | network  |
| CORE\_PEER\_GOSSIP\_CONNTIMEOUT                          | Connection timeout\.                                                                                                           | 0\.1\-2 s        | network  |
| CORE\_PEER\_GOSSIP\_RECVBUFFSIZE                         | Buffer size of received messages\.                                                                                             | 1\-20            | network  |
| CORE\_PEER\_GOSSIP\_SENDBUFFSIZE:                        | Buffer size of sending messages\.                                                                                              | 10\-200          | network  |
| CORE\_PEER\_GOSSIP\_DIGESTWAITTIME                       | Time to wait before pull engine processes incoming digests\.                                                                   | 0\.05\-1 s       | network  |
| CORE\_PEER\_GOSSIP\_REQUESTWAITTIME                      | Time to wait before pull engine removes incoming nonce\.                                                                       | 75\-1500 ms      | network  |
| CORE\_PEER\_GOSSIP\_RESPONSEWAITTIME                     | Time to wait before pull engine ends pull\.                                                                                    | 0\.1\-2 s        | network  |
| CORE\_PEER\_GOSSIP\_ALIVETIMEINTERVAL                    | Alive check interval\.                                                                                                         | 0\.25\-5 s       | network  |
| CORE\_PEER\_GOSSIP\_ALIVEEXPIRATIONTIMEOUT               | Alive expiration timeout\.                                                                                                     | 1\.25\-25 s      | network  |
| CORE\_PEER\_GOSSIP\_RECONNECTINTERVAL                    | Reconnect interval\.                                                                                                           | 1\.25\-25 s      | network  |
| CORE\_PEER\_GOSSIP\_ELECTION\_STARTUPGRACEPERIOD       | Longest time peer waits for stable membership during leader election startup                                                   | 0\.75\-15 s      | network  |
| CORE\_PEER\_GOSSIP\_ELECTION\_MEMBERSHIPSAMPLEINTERVAL | Interval gossip membership samples to check its stability                                                                      | 0\.05\-1 s       | network  |
| CORE\_PEER\_GOSSIP\_ELECTION\_LEADERALIVETHRESHOLD     | Time passes since last declaration message before peer decides to perform leader election                                      | 0\.5\-10 s       | network  |
| CORE\_PEER\_GOSSIP\_STATE\_CHECKINTERVAL               | CheckInterval interval to check whether peer is lagging behind enough to request blocks via state transfer from another peer\. | 0\.5\-10 s       | network  |
| CORE\_PEER\_GOSSIP\_STATE\_RESPONSETIMEOUT             | ResponseTimeout amount of time to wait for state transfer response from other peers\.                                          | 0\.15\-3 s       | network  |
| CORE\_PEER\_GOSSIP\_STATE\_BATCHSIZE                   | The number of blocks to request via state transfer from another peer\.                                                         | 0\.5\-10 s       | network  |
| CORE\_PEER\_GOSSIP\_STATE\_BLOCKBUFFERSIZE             | BlockBufferSize reflect the maximum distance between lowest and highest block sequence number state buffer to avoid holes\.    | 5\-100           | network  |
| CORE\_PEER\_GOSSIP\_STATE\_MAXRETRIES                  | MaxRetries maximum number of re\-tries to ask for single state transfer request                                                | 0\.15\-3         | network  |
|                                                                |                                                                                                                                |                  |          |
