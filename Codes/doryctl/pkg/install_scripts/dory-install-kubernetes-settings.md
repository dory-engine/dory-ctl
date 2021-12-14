# after installation

## dory-core settings after installed

### finish {{ $.dory.gitRepo.type }} install

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.port }}
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

{{- if eq $.dory.gitRepo.type "gitea" }}
- 1. open gitea url finish gitea install, set admin username / password / email
- 2. login to gitea, open `{{ $.viewURL }}:{{ $.dory.gitRepo.port }}/user/settings/applications`, generate new token.
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
    
### set all kubernetes nodes to connect {{ $.imageRepo.type }}

- add following /etc/hosts record in all kubernetes nodes for {{ $.imageRepo.type }} domain name 

```shell script
vi /etc/hosts
{{ $.hostIP }}  {{ $.imageRepo.domainName }}
```

- copy {{ $.imageRepo.type }} certificates to all kubernetes nodes

```shell script
scp -r /etc/docker/certs.d root@${KUBERNETES_HOST}:/etc/docker/
```

### restart dory-core and dory-dashboard

```shell script
# restart dory-core and dory-dashboard
kubectl -n {{ $.dory.namespace }} delete pods dory-core-0 dory-dashboard-0
```

## connect your dory

### dory-dashboard

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.port }}
- user: {{ $.dorycore.adminUser.username }}
- password file located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/dory-data/admin.password`
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/dory-core`

### {{ $.dory.gitRepo.type }}

- url: {{ $.viewURL }}:{{ $.dory.gitRepo.port }}
- data located at: `{{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.gitRepo.type }}`

### {{ $.imageRepo.type }}

- url: https://{{ $.imageRepo.domainName }}
- user: admin / {{ $.imageRepo.password }} (admin user)
- data located at: `{{ $.rootDir }}/{{ $.imageRepo.namespace }}`

### {{ $.dory.artifactRepo.type }}

- url: {{ $.viewURL }}:{{ $.dory.artifactRepo.port }}
- user: admin / Nexus_Pwd_321 (admin user)
- user: public-user / public-user (public user)

## about dory install files

- all dory kubernetes install files located at: `dory-install-kubernetes`