# doryctl 是一个安装和管理 `Dory-Engine` 的命令行工具

![](docs/images/dory-icon.png)

官方网站: [https://doryengine.com](https://doryengine.com)

## 什么是`Dory-Engine`

- `Dory-Engine` 是一个极简的应用上云引擎，应用开发者无需掌握复杂的DevOps和云原生知识，即可实现应用从源代码交付到云原生环境。 

## 使用 doryctl 安装 `Dory-Engine`

- Now we can use doryctl to install `Dory-Engine` with `docker-compose` (for test usage) or `kubernetes`(for production usage, recommended)
- 可以使用 doryctl 以两种方式安装 `Dory-Engine`。 
    1. 使用`docker-compose`把`Dory-Engine`安装在`docker`容器中，用于测试用途。
    2. 把`Dory-Engine`安装在`kubernetes`中，用于正式生产用途。

```shell script
##############################
# 根据以下指引把dory-core安装在docker中

# 1. 检查docker方式安装的环境是否就绪
doryctl install check --mode docker

# 2. 从docker hub拉取相关docker镜像
doryctl install pull

# 3. 打印docker方式安装的安装配置样例
doryctl install print --mode docker > install-config-docker.yaml

# 4. 根据环境的实际情况，修改安装配置
vi install-config-docker.yaml

# 5. 正式运行自动安装程序
doryctl install run -o readme-install-docker -f install-config-docker.yaml

##############################
# 根据以下指引把dory-core安装在kubernetes中

# 1. 检查kubernetes方式安装的环境是否就绪
doryctl install check --mode kubernetes

# 2. 从docker hub拉取相关docker镜像
doryctl install pull

# 3. 打印kubernetes方式安装的安装配置样例
doryctl install print --mode kubernetes > install-config-kubernetes.yaml

# 4. 根据环境的实际情况，修改安装配置
vi install-config-kubernetes.yaml

# 5. 正式运行自动安装程序
doryctl install run -o readme-install-kubernetes -f install-config-kubernetes.yaml
```

- 获取更多帮助请运行以下命令

```shell script
doryctl -h
```
