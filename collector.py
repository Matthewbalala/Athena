import requests
from utils.metrics import Metrics
import pandas as pd
import numpy as np
import pdb

class Collector(object):
    """
    收集caliper和prometheus信息
    params:
    endpoints: 必须包含peer和orderer两种类型的节点信息
    """

    def __init__(self, endpoints, reportfile, channel, chaincode):
        _nPeers  = 0
        _nOrderers = 0
        for item in endpoints:
            if item["node_type"] == "peer":
                _nPeers += 1
            elif item["node_type"] == "orderer":
                _nOrderers += 1
        assert(_nOrderers > 0)
        assert(_nPeers > 0)

        self._report = reportfile
        self._channel = channel
        self._chaincode = chaincode
        self._endpoints = endpoints
        self._filters = {
            "peer": [
                "chaincode_shim_request_duration_sum", 
                "chaincode_shim_request_duration_count",
                "deliver_blocks_sent",
                "endorser_proposal_duration_count",
                "endorser_proposal_duration_sum",
                "endorser_successful_proposals",
                # grpc + gossip
                "ledger_block_processing_time_sum",
                "ledger_block_processing_time_count",
                "ledger_blockchain_height",
                "ledger_blockstorage_and_pvtdata_commit_time_sum",
                "ledger_blockstorage_and_pvtdata_commit_time_count",
                "ledger_blockstorage_commit_time_sum",
                "ledger_blockstorage_commit_time_count",
                "ledger_statedb_commit_time_sum",
                "ledger_statedb_commit_time_count",
                "ledger_transaction_count",
                "logging_entries_checked",
                "logging_entries_written",
                "process_cpu_seconds_total",
                "process_max_fds",
                "process_open_fds",
                "process_resident_memory_bytes",
                "process_start_time_seconds",
                "process_virtual_memory_bytes",
                "process_virtual_memory_max_bytes",
            ],
            "orderer": [
                # "broadcast_enqueue_duration_sum",
                # "broadcast_enqueue_duration_count",
                # "broadcast_validate_duration_sum",
                # "broadcast_validate_duration_count",
                "cluster_comm_msg_send_time_sum",
                "cluster_comm_msg_send_time_count",
                "consensus_etcdraft_committed_block_number",
                "consensus_etcdraft_data_persist_duration_sum",
                "consensus_etcdraft_data_persist_duration_count",
                "consensus_etcdraft_is_leader",
                "consensus_etcdraft_leader_changes",
                "deliver_blocks_sent",
                "etcd_disk_wal_fsync_duration_seconds_sum",
                "etcd_disk_wal_fsync_duration_seconds_count",
                "ledger_blockstorage_commit_time_sum",
                "ledger_blockstorage_commit_time_count",
                "process_cpu_seconds_total",
                "process_max_fds",
                "process_resident_memory_byte",
            ],
            "peer-net": [
                "gossip_comm_messages_received",
                "gossip_comm_messages_sent",
                "gossip_membership_total_peers_known",
                "gossip_payload_buffer_size",
                "gossip_privdata_commit_block_duration_sum",
                "gossip_privdata_commit_block_duration_count",
                "gossip_privdata_fetch_duration_sum",
                "gossip_privdata_fetch_duration_count",
                "gossip_privdata_list_missing_duration_sum",
                "gossip_privdata_list_missing_duration_count",
                "gossip_privdata_purge_duration_sum",
                "gossip_privdata_purge_duration_count",
                "gossip_privdata_reconciliation_duration_sum",
                "gossip_privdata_reconciliation_duration_count",
                "gossip_privdata_validation_duration_sum",
                "gossip_privdata_validation_duration_count",
                "gossip_state_commit_duration_sum",
                "gossip_state_commit_duration_count",
                "gossip_state_height",
                "grpc_comm_conn_closed",
                "grpc_comm_conn_opened",
                "grpc_server_stream_messages_received",
                "grpc_server_stream_messages_sent",
                "grpc_server_stream_request_duration_sum",
                "grpc_server_stream_request_duration_count",
                "grpc_server_stream_requests_completed",
                "grpc_server_stream_requests_received",
                "grpc_server_unary_request_duration_sum",
                "grpc_server_unary_request_duration_count",
                "grpc_server_unary_requests_completed",
                "grpc_server_unary_requests_received",
            ],
            # "orderer-net": [
            #     "grpc_comm_conn_closed",
            #     "grpc_comm_conn_opened",
            #     "grpc_server_stream_messages_received",
            #     "grpc_server_stream_messages_sent",
            #     "grpc_server_stream_request_duration_sum",
            #     "grpc_server_stream_request_duration_count",
            # ]
        }

    def _convert_dict2name(self, input):
        """
        将metrics名称字典转换为_名称
        """
        res = ""
        for index in input:
            res = res + "_" + index + "_" + str(input[index])
        return res

    def collect_from_prometheus(self, handler=None):
        results = {
            "peer": [],
            "orderer": [],
            "peer-net":[],
            # "orderer-net":[]
        }

        _metricsobjs = self._get_metadata_from_prom()
        for node_type in results:
            for  metrics_interpreter in _metricsobjs[node_type]:
                metrics_lists = metrics_interpreter.interprete()
                # 数据清洗
                # pdb.set_trace()
                data = {}
                for item in metrics_lists:
                    if item['label'] == {}:
                        data[item['key']] = item['value']
                    elif 'chaincode' in item['label']:
                        # TODO: 此处限制使用智能合约
                        if item['label']['chaincode'].find(self._chaincode) != -1:
                            data[item['key'] + self._convert_dict2name(item['label'])] = item['value']

                    elif 'channel' in item['label']:
                        # TODO: 此处限制使用channel
                        if item['label']['channel'].find(self._channel) != -1:
                            data[item['key'] + '*' + self._convert_dict2name(item['label'])] = item['value']
                    else:
                        data[item['key']] = item['value']
                results[node_type].append(data)
        if handler:
            results = handler(results)
        return results

    def _get_metadata_from_prom(self):
        metrics_res = {
            "peer": [],
            "orderer": [],
            "peer-net":[],
            # "orderer-net":[]
        }
        for item in self._endpoints:
            if item["node_type"] == "peer":
                metrics_res["peer"].append(Metrics(item["url"], self._filters["peer"]))
                metrics_res["peer-net"].append(Metrics(item["url"], self._filters["peer-net"]))
            elif item["node_type"] == "orderer":
                metrics_res["orderer"].append(Metrics(item["url"], self._filters["orderer"]))
                # metrics_res["orderer-net"].append(Metrics(item["url"], self._filters["orderer-net"]))
        return metrics_res

    def collect_from_caliper(self):
        """
        针对operations提取metrics:
        1. TPS
        2. CPU(%)
        3. Memory(MB)
        4. Latency(s)
        """
        html = pd.read_html(self._report)
        orderer_flag = np.bitwise_not(html[2]["Name"].str.contains("orderer"))
        ca_flag = np.bitwise_not(html[2]["Name"].str.contains("ca"))
        data = html[2][orderer_flag & ca_flag].mean()
        CPU = data["CPU%(avg)"]
        Mem = data["Memory(avg) [MB]"]
        Latency = html[1]["Avg Latency (s)"][0]
        TPS = html[1]["Throughput (TPS)"][0]
        
        return {
            "CPU": CPU,
            "Mem": Mem,
            "Latency": Latency,
            "TPS": TPS
        }

    


# test
if __name__ == '__main__':
    urls = [{"node_type": "peer", "url":"http://10.10.7.46:7062/metrics"}, {"node_type": "orderer", "url":"http://10.10.7.52:7051/metrics"}]
    reportfile = 'caliper-deploy-tool/report.html'
    collector = Collector(urls, reportfile, "mychannel", "smallbank")
    print(collector._get_metadata_from_prom())
    print(len(collector.collect_from_prometheus()))
    print(collector.collect_from_prometheus())
    print(collector.collect_from_caliper())