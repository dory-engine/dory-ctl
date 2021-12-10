# doryctl

## doryctl install

- docker-compose模式安装

```shell script
# 检查节点上相关软件以及kubernetes集群是否可用
doryctl install check --mode docker

# 拉取相关镜像
doryctl install pull

# 获取安装设置
doryctl install print --mode docker
```
