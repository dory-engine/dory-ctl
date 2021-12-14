# doryctl

## 需求

- doryctl install reset 清空
- doryctl install run 显示后续操作提示，打印账号信息
- doryctl install check 显示创建k8s admin-user，以及获取token命令
- doryctl install manual 创建手工安装文件以及说明文档，可以自行修改进行定制化部署

- kuberentes gitea启动异常（权限异常）
- kubernetes docker的harbor证书异常

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
vi docker_install.yaml

# 执行安装
doryctl install run --mode docker -f docker_install.yaml
```

1. 打开gitea界面，设置gitea管理员账号密码，完成gitea安装。然后获取gitea的管理员token。
2. 使用gitea管理员的账号密码以及token信息，设置dory/dory-core/config/config.yaml中`PLEASE_INPUT_BY_MANUAL`设置项。
3. 把/etc/docker/certs.d/中对应的harbor证书复制到所有k8s集群节点。同时必须保证所有k8s集群节点包含/etc/timezone和/etc/localtime时区设置文件。
4. 保证所有k8s节点的/etc/hosts包含harbor域名
5. 重启dory-core。docker rm -f dory-core && docker-compose up -d

### kubernetes模式安装

```shell script
# 创建k8s admin-user serviceAccount

# 创建证书
mkdir -p /etc/docker/certs.d/registry2.new.imdory.com/
kubectl -n harbor get secrets harbor-ingress -o jsonpath='{ .data.ca\.crt }' | base64 -d > /etc/docker/certs.d/registry2.new.imdory.com/ca.crt
kubectl -n harbor get secrets harbor-ingress -o jsonpath='{ .data.tls\.crt }' | base64 -d > /etc/docker/certs.d/registry2.new.imdory.com/registry2.new.imdory.com.cert
kubectl -n harbor get secrets harbor-ingress -o jsonpath='{ .data.tls\.key }' | base64 -d > /etc/docker/certs.d/registry2.new.imdory.com/registry2.new.imdory.com.key

```
