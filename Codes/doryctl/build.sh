# 执行编译
cd /root/devops-dory-ctl/Codes/doryctl
git pull
date && time go build
ls -alh

# # 复制所有到new1-dory
# scp -r doryctl root@new1-dory:/usr/bin/
# scp -r pkg/install_configs/*.yaml root@new1-dory:/root

# # 复制所有到new2-dory
# scp -r doryctl root@new2-dory:/usr/bin/
# scp -r pkg/install_configs/*.yaml root@new2-dory:/root
