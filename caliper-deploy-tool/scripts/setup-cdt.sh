echo "Set variable CDTHOME: ${CDTHOME}"

docker pull node:10-alpine
docker pull coredns/coredns:1.8.0
docker pull itsthenetwork/nfs-server-alpine:12


echo "1. build cdt docker image."
docker build -t cdt ${CDTHOME}/docker

echo "2. stop systemd-resolved and release dns port 53"
service systemd-resolved stop

echo "3. boot coredns server."
docker run -it -d --net=host \
--name=coredns --restart=always \
-v ${CDTHOME}/dns/coredns:/etc/coredns/ \
coredns/coredns:1.8.0 \
-conf /etc/coredns/Corefile

echo "4. boot nfs server."
docker run -d --name nfs \
--privileged \
-v ${CDTHOME}/benchmarks/config_raft:/nfsshare \
-e SHARED_DIRECTORY=/nfsshare \
-p 2049:2049 itsthenetwork/nfs-server-alpine:12

echo "Finish."
