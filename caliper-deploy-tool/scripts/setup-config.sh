cd ansible
if [ ${CDTIP} ];then
	echo "Set NFS's ip: ${CDTIP}"
else
	echo -e "\033[31m CDTIP IS NOT SET \033[0m"
    exit
fi


echo "generate mount-nfs.yaml"
sed "s/ip/${CDTIP}/g" mount-nfs.tpl > mount-nfs.yaml

echo "update docker config"
ansible-playbook ${CDTHOME}/ansible/update-docker.yaml

echo "apt install nfs-common"
echo "mount nfs to /root/ansible/nfs"
ansible-playbook ${CDTHOME}/ansible/umount-nfs.yaml
ansible-playbook ${CDTHOME}/ansible/mount-nfs.yaml