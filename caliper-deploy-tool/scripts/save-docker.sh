export ORDERERTAG=1.4
export PEERTAG=1.4
export CATAG=1.4
export CCENVTAG=1.4
export TARGETTAG=cdt
export DISTDIR=/root/ansible/images

rm -rf $DISTDIR/*

mkdir -p $DISTDIR

echo "Pull image from dockerhub..."
docker pull hyperledger/fabric-peer:$PEERTAG
docker pull hyperledger/fabric-ca:$CATAG
docker pull hyperledger/fabric-orderer:$ORDERERTAG
docker pull hyperledger/fabric-ccenv:$CCENVTAG
echo "Pull finished. Save and tag docker images."

docker tag hyperledger/fabric-orderer:$ORDERERTAG hyperledger/fabric-orderer:$TARGETTAG
docker tag hyperledger/fabric-peer:$PEERTAG hyperledger/fabric-peer:$TARGETTAG
docker tag hyperledger/fabric-ca:$CATAG hyperledger/fabric-ca:$TARGETTAG
docker tag hyperledger/fabric-ccenv:$CCENVTAG hyperledger/fabric-ccenv:latest


docker save hyperledger/fabric-peer:$TARGETTAG -o $DISTDIR/peer-docker.tar
docker save hyperledger/fabric-orderer:$TARGETTAG -o $DISTDIR/orderer-docker.tar
docker save hyperledger/fabric-ca:$TARGETTAG -o $DISTDIR/ca-docker.tar
docker save hyperledger/fabric-ccenv:latest -o $DISTDIR/ccenv-docker.tar