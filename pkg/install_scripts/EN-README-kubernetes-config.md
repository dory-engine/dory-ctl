# config dory after install in kubernetes

## dory-core settings after installed

{{- if $.dory.gitRepo.internal.image }}
### finish {{ $.dory.gitRepo.type }} install and update dory config.yaml

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

{{- if eq $.dory.gitRepo.type "gitea" }}
- 1. open gitea url finish gitea install, at `Administrator Account Settings ` set admin username / password / email
- 2. login to gitea, open `{{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}/user/settings/applications`, at `Generate New Token` generate a new token.
{{- else if eq $.dory.gitRepo.type "gitlab" }}
- 1. gitlab password file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}/config/initial_root_password`
- 2. login to gitlab, open `{{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}/-/profile/personal_access_tokens`, add a personal access token.
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
{{- end }}

{{- if $.artifactRepoInternal }}
### update {{ $.dory.artifactRepo.type }} admin password and update dory config.yaml

- url: {{ $.viewURL }}:{{ $.dory.artifactRepo.internal.port }}
- user: admin / {{ $.artifactRepoPassword }} (admin user)
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.artifactRepo.type }}`

- 1. open {{ $.dory.artifactRepo.type }} url, login as admin user
- 2. change admin password, open `{{ $.viewURL }}:{{ $.dory.artifactRepo.internal.port }}/#user/account` and change password
- 3. update dory-core config file:
  - config file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/config/config.yaml`
  - search `{{ $.artifactRepoPassword }}` in config file
  - update following admin user password: 
    - artifactRepoConfigs.password
{{- end }}
 
### set all kubernetes nodes to connect {{ $.imageRepo.type }}

- 1. add following {{ $.imageRepo.type }} domain name in /etc/hosts record for all kubernetes nodes  

```shell script
vi /etc/hosts
{{ $.imageRepoIp }}  {{ $.imageRepoDomainName }}
```

- 2. copy {{ $.imageRepo.type }} certificates to all kubernetes nodes

```shell script
# copy certificates to all kubernetes nodes
scp -r /etc/docker/certs.d root@${KUBERNETES_HOST}:/etc/docker/

{{- if not $.imageRepoInternal }}
# copy {{ $.imageRepo.type }} certificates to docker executor directory
cp -rp /etc/docker/certs.d/{{ $.imageRepoDomainName }} {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}
{{- end }}
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

{{- if $.dory.gitRepo.internal.image }}
### {{ $.dory.gitRepo.type }} git repository

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.internal.port }}
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`
{{- end }}

{{- if $.artifactRepoInternal }}
### {{ $.dory.artifactRepo.type }} artifact and dependency repository

- url: {{ $.viewURL }}:{{ $.dory.artifactRepo.internal.port }}
- public user: {{ $.artifactRepoPublicUser }} / {{ $.artifactRepoPublicPassword }}
- docker.io image proxy: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortHub }}
- gcr.io image proxy: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortGcr }}
- quay.io image proxy: {{ $.artifactRepoIp }}:{{ $.artifactRepoPortQuay }}
{{- end }}

{{- if $.imageRepoInternal }}
### {{ $.imageRepo.type }} image repository

- url: https://{{ $.imageRepoDomainName }}
- user: admin / {{ $.imageRepoPassword }} (admin user)
- data located at: `{{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}`
{{- end }}

### openldap account management

- url: {{ $.viewURL | replace "http://" "https://" }}:{{ $.dory.openldap.port }}
- user: cn=admin,{{ $.dory.openldap.baseDN }} / {{ $.dory.openldap.password }}

### caution: this folder is very important, included all config files and readme files, please keep it
