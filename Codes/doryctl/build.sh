# 执行编译
cd /root/devops-dory-core/Codes/doryctl
git pull
date && time go build
ls -al

# # 复制所有到new1-dory
# scp -r doryctl root@new1-dory:/usr/bin/

# # 复制所有到new2-dory
# scp -r doryctl root@new2-dory:/usr/bin/
