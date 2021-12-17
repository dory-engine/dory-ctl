# 执行编译
cd /root/devops-dory-ctl/Codes/doryctl
git pull
date && time go build -o doryctl
ls -alh

# # 复制所有到new1-dory
# scp -r doryctl root@new1-dory:/usr/bin/

# # 复制所有到new2-dory
# scp -r doryctl root@new2-dory:/usr/bin/
