# 以docker方式部署完dory之后，必须进行以下设置

## 安装完成后必须进行dory-core配置

### 在kubernetes的共享存储中创建目录

- 创建目录

{{- if $.kubernetes.pvConfigLocal.localPath }}
```shell script
# 在kubernetes的本地存储创建目录
mkdir -p {{ $.kubernetes.pvConfigLocal.localPath }}
```
{{- else if $.kubernetes.pvConfigNfs.nfsPath }}
```shell script
# 在kubernetes的nfs存储创建目录
mkdir -p {{ $.kubernetes.pvConfigNfs.nfsPath }}
```
{{- else if $.kubernetes.pvConfigCephfs.cephPath }}
```shell script
# 在kubernetes的cephfs存储创建目录
mkdir -p {{ $.kubernetes.pvConfigCephfs.cephPath }}
```
{{- end }}

- 重启project-data-alpine-0

```shell script
kubectl -n {{ $.dory.namespace }} delete pods project-data-alpine-0

# 检查project-data-alpine-0是否正常
kubectl -n {{ $.dory.namespace }} get pods project-data-alpine-0
```

### 完成 {{ $.dory.gitRepo.type }} 安装并更新dory的config.yaml配置

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.port }}
- 数据存放在以下目录: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

{{- if eq $.dory.gitRepo.type "gitea" }}
- 1. 打开gitea的网址，完成gitea安装设置，重点设置 `管理员账号` ，设置管理员的用户名、密码、邮箱
- 2. 登录gitea，打开 `{{ $.viewURL }}:{{ $.dory.gitRepo.port }}/user/settings/applications`，在`创建新Token`创建一个新的访问token。
{{- else if eq $.dory.gitRepo.type "gitlab" }}
- 1. gitlab的root用户密码文件存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}/config/initial_root_password`
- 2. 登录gitlab，打开 `{{ $.viewURL }}:{{ $.dory.gitRepo.port }}/-/profile/personal_access_tokens`，新增一个访问token。
{{- end }}
- 3. 记住管理员的 `用户名、密码、邮箱、访问token` 用于更新dory-core的配置文件中的 {{ $.dory.gitRepo.type }} 设置
- 4. 更新dory-core配置文件:
  - 配置文件存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/config/config.yaml`
  - 配置文件中搜索: `PLEASE_INPUT_BY_MANUAL`
  - 更新配置文件中以下代码仓库管理员设置: 
    - gitRepoConfigs.username
    - gitRepoConfigs.name
    - gitRepoConfigs.mail
    - gitRepoConfigs.password
    - gitRepoConfigs.token
    
### 更新 {{ $.dory.artifactRepo.type }} 管理员密码，并更新dory的config.yaml配置文件

- url: {{ $.viewURL }}:{{ $.dory.artifactRepo.port }}
- user: admin / Nexus_Pwd_321 (管理员用户)
- 数据存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.artifactRepo.type }}`

- 1. 打开 {{ $.dory.artifactRepo.type }} 网址，使用admin的默认账号密码登录
- 2. 修改管理员密码: `{{ $.viewURL }}:{{ $.dory.artifactRepo.port }}/#user/account`
- 3. 更新dory-core配置文件:
  - 配置文件存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/config/config.yaml`
  - 在配置文件中搜索 `Nexus_Pwd_321`
  - 更新以下管理员密码配置: 
    - artifactRepoConfigs.password
 
### 设置所有kubernetes节点连接 {{ $.imageRepo.type }}

- 1. 添加以下 {{ $.imageRepo.type }} 域名记录到所有kubernetes节点的 /etc/hosts 文件  

```shell script
vi /etc/hosts
{{ $.hostIP }}  {{ $.imageRepo.domainName }}
```

- 2. 复制 {{ $.imageRepo.type }} 证书到所有kubernetes节点

```shell script
scp -r /etc/docker/certs.d root@${KUBERNETES_HOST}:/etc/docker/
```

### 重启 dory-core 和 dory-dashboard 服务

- 1. 重启 dory-core 和 dory-dashboard 服务

```shell script
cd {{ $.rootDir }}/{{ $.dory.namespace }}
docker rm -f dory-core && docker-compose up -d
```

## 访问各个dory服务

### dory-dashboard

- url: {{ $.viewURL }}:{{ $.dory.dorycore.port }}
- 管理员用户: {{ $.dorycore.adminUser.username }}
- 管理员账号密码存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/dory-data/admin.password`
- dory-core数据和配置存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core`

### {{ $.dory.gitRepo.type }}

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.port }}
- 数据存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

### {{ $.dory.artifactRepo.type }}

- url: {{ $.viewURL }}:{{ $.dory.artifactRepo.port }}
- user: public-user / public-user (公共用户账号)

### {{ $.imageRepo.type }}

- url: https://{{ $.imageRepo.domainName }}
- user: admin / {{ $.imageRepo.password }} (管理员用户)
- 数据存放在: `{{ $.rootDir }}/{{ $.imageRepo.namespace }}`

### 注意，本目录非常重要，本目录为安装过程配置文件以及说明文件目录，建议保留
