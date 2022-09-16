export CDTHOME=${PWD}

setup-cdt:
	cd scripts && bash setup-cdt.sh

purge:
	cd scripts && bash purge-cdt.sh

distribute-docker-images:
	@echo "save docker images to tar file"
	bash scripts/save-docker.sh
	@echo "distribute images to each host"
	ansible-playbook ansible/update-image.yaml

setup-config:
	bash scripts/setup-config.sh


generate:
	@echo "1. generate dist config"
	@echo "2. generate fabric crypto-config"
	python3 scripts/export-config.py
	mv dist/configtx.yaml benchmarks/config_raft
	mv dist/crypto-config.yaml benchmarks/config_raft
	cd benchmarks/config_raft && bash generate.sh

deploy-fabric-up:
	@echo "deploy & boot fabric network."
	ansible-playbook ansible/deploy-up.yaml

deploy-fabric-down:
	@echo "stop fabric network & clean all docker containers"
	ansible-playbook ansible/deploy-down.yaml

start-cdt:
	@echo "booting caliper using cdt to test.."
	bash scripts/boot.sh


	


