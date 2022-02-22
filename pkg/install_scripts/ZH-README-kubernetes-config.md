# 以kubernetes方式部署完dory之后，必须进行以下设置

## 安装完成后必须进行dory-core配置

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
kubectl -n {{ $.dory.namespace }} delete pods dory-core-0 dory-dashboard-0

# 等待 dory-core-0 dory-dashboard-0 pod处于ready状态
kubectl -n {{ $.dory.namespace }} get pods -o wide -w
```

## 访问各个dory服务

### dory-dashboard 管理界面

- url: {{ $.viewURL }}:{{ $.dory.dorycore.port }}
- 管理员用户: {{ $.dorycore.adminUser.username }}
- 管理员账号密码存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/dory-data/admin.password`
- dory-core数据和配置存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core`

### {{ $.dory.gitRepo.type }} 代码仓库

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.port }}
- 数据存放在: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

### {{ $.dory.artifactRepo.type }} 依赖与制品仓库

- url: {{ $.viewURL }}:{{ $.dory.artifactRepo.port }}
- 公共用户账号: public-user / public-user
- docker.io镜像代理地址: {{ $.hostIP }}:{{ $.dory.artifactRepo.portHub }}
- gcr.io镜像代理地址: {{ $.hostIP }}:{{ $.dory.artifactRepo.portGcr }}
- quay.io镜像代理地址: {{ $.hostIP }}:{{ $.dory.artifactRepo.portQuay }}

### {{ $.imageRepo.type }} 容器镜像仓库

- url: https://{{ $.imageRepo.domainName }}
- user: admin / {{ $.imageRepo.password }} (管理员用户)
- 数据存放在: `{{ $.rootDir }}/{{ $.imageRepo.namespace }}`

### openldap 账号管理中心

- url: {{ $.viewURL | replace "http://" "https://" }}:{{ $.dory.openldap.port }}
- 管理员用户: cn=admin,{{ $.dory.openldap.baseDN }} / {{ $.dory.openldap.password }}

### 注意，本目录非常重要，本目录为安装过程配置文件以及说明文件目录，建议保留
