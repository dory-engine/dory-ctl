# Dory-Engine是一个简单得非常吓人的应用上云引擎，doryctl是Dory-Engine的命令行管理工具

- [English documents](README.md)
- [中文文档](README-zh.md)

- [Github](https://github.com/dory-engine/dory-dashboard)
- [Gitee](https://gitee.com/dory-engine/dory-dashboard)

![](docs/images/dory-icon.png)

详细参见官方网站: [https://doryengine.com](https://doryengine.com)

## 什么是`Dory-Engine`

![](docs/images/what-is-dory.png)

- `Dory-Engine` 是一个极简的应用上云引擎

- 应用开发者无需掌握复杂的DevOps和云原生知识，即可实现应用从源代码交付到云原生环境。

### `Dory-Engine`架构

![](docs/images/architecture.png)

- 分布式: Dory-Engine使用无状态设计架构，可部署在Kubernetes或者docker中，轻松实现分布式水平扩缩容。
- 全容器: 步骤都在远程步骤执行器(Docker)中执行，可以轻松实现负载分担。
- 高弹性: 远程步骤执行器(Docker)可以根据工作负载，进行水平扩缩容实现高弹性。
- 易扩展: 通过容器技术，让步骤支持各种执行环境，实现应用上云流程的灵活扩展。
- 多云编排: 可以同时接管多个不同的云原生环境、主机环境(企业版原生支持)、各种数据库环境(企业版原生支持)，把应用发布到多个不同环境。
- 协同治理: 接管DevOps持续交付工具链各个组件，自动开通配置好各个组件和云原生环境，应用上云从未如此简单。

## 安装 doryctl

### 二进制包方式安装

- 从[doryctl Releases](https://github.com/dory-engine/dory-ctl/releases)下载相应的二进制安装包. 一旦下载后，二进制执行文件就可以在任何地方执行，你无需把执行文件放到全局目录，假如你没有主机的相关权限的话。

- 更好的方式是把二进制执行文件放到PATH所在的目录，这样你在任何地方都可以调用doryctl，常见的路径为 /user/local/bin

### 源码方式安装

- doryctl 需要golangv1.16.x以上版本编译。

```shell script
# 克隆源代码
git clone https://github.com/dory-engine/dory-ctl.git

# 从源代码编译doryctl
cd dory-ctl
go mod tidy && go build -o doryctl
mv doryctl /usr/local/bin/
```


## 使用 doryctl 安装 `Dory-Engine`

- 可以使用 doryctl 以两种方式安装 `Dory-Engine`: 
    - 1. 使用`docker-compose`把`Dory-Engine`安装在`docker`容器中，用于测试用途。
    - 2. 把`Dory-Engine`安装在`kubernetes`中，用于正式生产用途。

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
