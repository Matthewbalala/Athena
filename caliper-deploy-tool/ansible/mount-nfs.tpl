---
- hosts: cdt
  remote_user: root
  vars:
    dest_dir: /root/ansible/nfs
    nfs: ip
  tasks:
  - name: mkdir into {{dest_dir}}
    shell: mkdir -p {{dest_dir}}
    ignore_errors: True

  - name: apt install nfs-common
    shell: apt install -y nfs-common 
    ignore_errors: True

  - name: mount {{nfs}}:/ to  {{dest_dir}}
    shell: mount -t nfs {{nfs}}:/ {{dest_dir}}