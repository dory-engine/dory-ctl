# doryctl

## doryctl install 安装dory-core组件

### docker-compose模式安装

```shell script
# 检查节点上相关软件以及kubernetes集群是否可用
doryctl install check --mode docker

# 拉取相关镜像
doryctl install pull

# 获取安装设置
doryctl install print --mode docker > docker_install.yaml

# 根据实际情况修改安装设置

# 执行安装
doryctl install run --mode docker -f docker_install.yaml
```

1. 打开gitea界面，设置gitea管理员账号密码，完成gitea安装。然后获取gitea的管理员token。
2. 使用gitea管理员的账号密码以及token信息，设置dory/dory-core/config/config.yaml中`PLEASE_INPUT_BY_MANUAL`设置项。
3. 把/etc/docker/certs.d/中对应的harbor证书复制到所有k8s集群节点。同时必须保证所有k8s集群节点包含/etc/timezone和/etc/localtime时区设置文件。
4. 在目标k8s集群部署project-data-alpine statefulset。
