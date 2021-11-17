echo "Set variable CDTHOME: ${CDTHOME}"

echo "1. rm cdt docker image."
docker rmi cdt

echo "2. stop & rm coredns server."
docker stop coredns && docker rm coredns

echo "3. umount nfs"
ansible-playbook ${CDTHOME}/ansible/umount-nfs.yaml

echo "4. stop & rm nfs server."
docker stop nfs && docker rm nfs

echo "5. rm -rf /root/ansible  on each host"
echo "6. stop&rm fabric's docker container on each host"
# 7. (non-implements) restore docker config on each host
rm -rf /root/ansible
ansible-playbook ${CDTHOME}/ansible/clean-file.yaml
service systemd-resolved start


echo "Finish."
