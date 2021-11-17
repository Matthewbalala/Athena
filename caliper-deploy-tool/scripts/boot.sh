echo "Set variable CDTHOME: ${CDTHOME}"

if [ ${CDTIP} ];then
	echo "Set DNS's ip: ${CDTIP}"
else
	echo -e "\033[31m CDTIP IS NOT SET \033[0m"
    exit
fi


docker rm cdt

docker run -v  ${CDTHOME}:/hyperledger/caliper/workspace \
--dns ${CDTIP} --name=cdt \
cdt npx caliper launch master \
--caliper-workspace  /hyperledger/caliper/workspace \
--caliper-benchconfig dist/config-distributed.yaml \
--caliper-networkconfig dist/client.yaml


echo "Test finish and benchmark result will be saved into ${CDTHOME}/report.html."