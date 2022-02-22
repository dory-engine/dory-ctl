# config dory after install in kubernetes

## dory-core settings after installed

### finish {{ $.dory.gitRepo.type }} install and update dory config.yaml

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.port }}
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

{{- if eq $.dory.gitRepo.type "gitea" }}
- 1. open gitea url finish gitea install, at `Administrator Account Settings ` set admin username / password / email
- 2. login to gitea, open `{{ $.viewURL }}:{{ $.dory.gitRepo.port }}/user/settings/applications`, at `Generate New Token` generate a new token.
{{- else if eq $.dory.gitRepo.type "gitlab" }}
- 1. gitlab password file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}/config/initial_root_password`
- 2. login to gitlab, open `{{ $.viewURL }}:{{ $.dory.gitRepo.port }}/-/profile/personal_access_tokens`, add a personal access token.
{{- end }}
- 3. copy admin `username / password / email / token` to update dory-core config file {{ $.dory.gitRepo.type }} settings
- 4. update dory-core config file:
  - config file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/config/config.yaml`
  - search `PLEASE_INPUT_BY_MANUAL` in config file
  - update following admin user settings: 
    - gitRepoConfigs.username
    - gitRepoConfigs.name
    - gitRepoConfigs.mail
    - gitRepoConfigs.password
    - gitRepoConfigs.token
    
### update {{ $.dory.artifactRepo.type }} admin password and update dory config.yaml

- url: {{ $.viewURL }}:{{ $.dory.artifactRepo.port }}
- user: admin / Nexus_Pwd_321 (admin user)
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.artifactRepo.type }}`

- 1. open {{ $.dory.artifactRepo.type }} url, login as admin user
- 2. change admin password, open `{{ $.viewURL }}:{{ $.dory.artifactRepo.port }}/#user/account` and change password
- 3. update dory-core config file:
  - config file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/config/config.yaml`
  - search `Nexus_Pwd_321` in config file
  - update following admin user password: 
    - artifactRepoConfigs.password
 
### set all kubernetes nodes to connect {{ $.imageRepo.type }}

- 1. add following {{ $.imageRepo.type }} domain name in /etc/hosts record for all kubernetes nodes  

```shell script
vi /etc/hosts
{{ $.hostIP }}  {{ $.imageRepo.domainName }}
```

- 2. copy {{ $.imageRepo.type }} certificates to all kubernetes nodes

```shell script
scp -r /etc/docker/certs.d root@${KUBERNETES_HOST}:/etc/docker/
```

### restart dory-core and dory-dashboard

- 1. restart dory-core and dory-dashboard

```shell script
kubectl -n {{ $.dory.namespace }} delete pods dory-core-0 dory-dashboard-0

# waiting for dory-core-0 dory-dashboard-0 ready
kubectl -n {{ $.dory.namespace }} get pods -o wide -w
```

## connect your dory

### dory-dashboard admin dashboard

- url: {{ $.viewURL }}:{{ $.dory.dorycore.port }}
- user: {{ $.dorycore.adminUser.username }}
- password file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/dory-data/admin.password`
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core`

### {{ $.dory.gitRepo.type }} git repository

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.port }}
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

### {{ $.dory.artifactRepo.type }} artifact and dependency repository

- url: {{ $.viewURL }}:{{ $.dory.artifactRepo.port }}
- public user: public-user / public-user
- docker.io image proxy: {{ $.hostIP }}:{{ $.dory.artifactRepo.portHub }}
- gcr.io image proxy: {{ $.hostIP }}:{{ $.dory.artifactRepo.portGcr }}
- quay.io image proxy: {{ $.hostIP }}:{{ $.dory.artifactRepo.portQuay }}

### {{ $.imageRepo.type }} image repository

- url: https://{{ $.imageRepo.domainName }}
- user: admin / {{ $.imageRepo.password }} (admin user)
- data located at: `{{ $.rootDir }}/{{ $.imageRepo.namespace }}`

### openldap account management

- url: {{ $.viewURL | replace "http://" "https://" }}:{{ $.dory.openldap.port }}
- user: cn=admin,{{ $.dory.openldap.baseDN }} / {{ $.dory.openldap.password }}

### caution: this folder is very important, included all config files and readme files, please keep it
