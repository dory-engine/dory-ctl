# 以docker-compose方式部署dory

## 概要

1. 请根据 `README-docker-install.md` 的说明手工安装dory
2. 请根据 `README-docker-config.md` 的说明在完成安装后手工设置dory
3. 假如安装失败，请根据 `README-docker-reset.md` 的说明停止所有dory服务并重新安装

## 复制所有脚本、配置文件到安装目录

```shell script
# 复制所有脚本、配置文件到安装目录
mkdir -p {{ $.rootDir }}
cp -rp * {{ $.rootDir }}
```

## 使用docker-compose部署 {{ $.imageRepo.type }}

```shell script
# 创建 {{ $.imageRepo.type }} 自签名证书
cd {{ $.rootDir }}/{{ $.imageRepo.namespace }}
sh harbor_certs.sh
ls -alh

# 安装 {{ $.imageRepo.type }}
cd {{ $.rootDir }}/{{ $.imageRepo.namespace }}
chmod a+x common.sh install.sh prepare
sh install.sh
ls -alh

# 停止并更新 {{ $.imageRepo.type }} 的 docker-compose.yml 文件
sleep 5 && docker-compose stop && docker-compose rm -f
export HARBOR_CONFIG_ROOT_PATH=$(echo "{{ $.rootDir }}/{{ $.imageRepo.namespace }}" | sed 's#\/#\\\/#g')
sed -i "s/${HARBOR_CONFIG_ROOT_PATH}/./g" docker-compose.yml
cat docker-compose.yml

# 重启 {{ $.imageRepo.type }}
docker-compose up -d
sleep 10

# 检查 {{ $.imageRepo.type }} 状态
docker-compose ps

# 在当前主机以及所有kubernetes节点主机上，把 {{ $.imageRepo.type }} 的域名记录添加到 /etc/hosts
vi /etc/hosts
{{ $.hostIP }}  {{ $.imageRepo.domainName }}

# 设置docker客户端登录到 {{ $.imageRepo.type }}
docker login --username admin --password {{ $.imageRepo.password }} {{ $.imageRepo.domainName }}

# 在 {{ $.imageRepo.type }} 中创建 public, hub, gcr, quay 四个项目
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "public", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "hub", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "gcr", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "quay", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'

# 把之前拉取的docker镜像推送到 {{ $.imageRepo.type }}
{{- range $_, $image := $.dockerImages }}
docker tag {{ if $image.dockerFile }}{{ $image.target }}{{ else }}{{ $image.source }}{{ end }} {{ $.imageRepo.domainName }}/{{ $image.target }}
{{- end }}

{{- range $_, $image := $.dockerImages }}
docker push {{ $.imageRepo.domainName }}/{{ $image.target }}
{{- end }}
```

## 使用docker-compose方式安装dory组件

```shell script
# 创建 docker executor 自签名证书
cd {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}
sh docker_certs.sh
ls -alh

# 从docker镜像中复制nexus初始化数据
cd {{ $.rootDir }}/{{ $.dory.namespace }}
docker rm -f nexus-data-init || true
docker run -d -t --name nexus-data-init doryengine/nexus-data-init:alpine-3.15.0 cat
docker cp nexus-data-init:/nexus-data/nexus .
docker rm -f nexus-data-init
chown -R 200:200 nexus
ls -alh nexus

# 创建 dory 组件目录并设置权限
cd {{ $.rootDir }}/{{ $.dory.namespace }}
mkdir -p mongo-core-dory
chown -R 999:999 mongo-core-dory
mkdir -p dory-core/dory-data
mkdir -p dory-core/tmp
chown -R 1000:1000 dory-core
ls -alh

# 使用docker-compose启动所有dory组件
cd {{ $.rootDir }}/{{ $.dory.namespace }}
ls -alh
docker-compose stop && docker-compose rm -f && docker-compose up -d

# 检查dory组件的状态
sleep 10
docker-compose ps
```

## 在kubernetes中创建project-data-alpine pod

```shell script
# project-data-alpine pod 用于创建项目的应用文件目录
# 在kubernetes中创建project-data-alpine pod
cd {{ $.rootDir }}
kubectl apply -f project-data-alpine.yaml
kubectl -n {{ $.dory.namespace }} get pods
```

## 请继续完成dory的配置

2. 请根据 `README-docker-config.md` 的说明在完成安装后手工设置dory
