# install dory with kubernetes

## summary

1. please follow `README-kubernetes-install.md` to install dory by manual
2. please follow `README-kubernetes-config.md` to config dory by manual after install
3. if install fail, please follow `README-kubernetes-reset.md` to stop all dory services and install again

## create install root directories

{{- $harborDomainName := $.imageRepo.internal.domainName }}
{{- $harborUserName := "admin" }}
{{- $harborPassword := $.imageRepo.internal.password }}
{{- if $.imageRepo.external.url }}{{ $harborDomainName = $.imageRepo.external.url }}{{ end }}
{{- if $.imageRepo.external.username }}{{ $harborUserName = $.imageRepo.external.username }}{{ end }}
{{- if $.imageRepo.external.password }}{{ $harborPassword = $.imageRepo.external.password }}{{ end }}

```shell script
{{- if $.imageRepo.internal.domainName }}
# create {{ $.imageRepo.type }} root directory
mkdir -p {{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}/database
mkdir -p {{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}/jobservice
mkdir -p {{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}/redis
mkdir -p {{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}/registry
chown -R 999:999 {{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}/database
chown -R 10000:10000 {{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}/jobservice
chown -R 999:999 {{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}/redis
chown -R 10000:10000 {{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}/registry
ls -alh {{ $.rootDir }}/{{ $.imageRepo.internal.namespace }}
{{- end }}

# create dory root directory
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/dory-data
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/tmp
cp -rp {{ $.dory.namespace }}/{{ $.dory.docker.dockerName }} {{ $.rootDir }}/{{ $.dory.namespace }}/
cp -rp {{ $.dory.namespace }}/dory-core {{ $.rootDir }}/{{ $.dory.namespace }}/
chown -R 1000:1000 {{ $.rootDir }}/{{ $.dory.namespace }}/dory-core
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/mongo-core-dory
chown -R 999:999 {{ $.rootDir }}/{{ $.dory.namespace }}/mongo-core-dory
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}

```

## {{ $.imageRepo.type }} installation and configuration

```shell script
{{- if $.imageRepo.internal.domainName }}
# create {{ $.imageRepo.type }} namespace and pv
kubectl delete ns {{ $.imageRepo.internal.namespace }}
kubectl delete pv {{ $.imageRepo.internal.namespace }}-pv
kubectl apply -f {{ $.imageRepo.internal.namespace }}/step01-namespace-pv.yaml

# install {{ $.imageRepo.type }}
helm install -n {{ $.imageRepo.internal.namespace }} {{ $.imageRepo.internal.namespace }} {{ $.imageRepo.type }}
helm -n {{ $.imageRepo.internal.namespace }} list

# waiting for all {{ $.imageRepo.type }} services ready
kubectl -n {{ $.imageRepo.internal.namespace }} get pods -o wide

# create {{ $.imageRepo.type }} self signed certificates and copy to /etc/docker/certs.d
sh {{ $.imageRepo.internal.namespace }}/harbor_update_docker_certs.sh
ls -alh /etc/docker/certs.d/{{ $.imageRepo.internal.domainName }}
{{- else }}
# copy harbor server ({{ $.imageRepo.external.ip }}) certificates to this node /etc/docker/certs.d/{{ $.imageRepo.external.url }} directory
# certificates are: ca.crt, {{ $.imageRepo.external.url }}.cert, {{ $.imageRepo.external.url }}.key
{{- end }}

# on current host and all kubernetes nodes add {{ $.imageRepo.type }} domain name in /etc/hosts
vi /etc/hosts
{{- if $.imageRepo.internal.domainName }}
{{ $.hostIP }}  {{ $harborDomainName }}
{{- else }}
{{ $.imageRepo.external.ip }}  {{ $harborDomainName }}
{{- end }}

# docker login to {{ $.imageRepo.type }}
docker login --username {{ $harborUserName }} --password {{ $harborPassword }} {{ $harborDomainName }}

# create public, hub, gcr, quay projects in {{ $.imageRepo.type }}
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "public", "public": true}' 'https://{{ $harborUserName }}:{{ $harborPassword }}@{{ $harborDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "hub", "public": true}' 'https://{{ $harborUserName }}:{{ $harborPassword }}@{{ $harborDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "gcr", "public": true}' 'https://{{ $harborUserName }}:{{ $harborPassword }}@{{ $harborDomainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "quay", "public": true}' 'https://{{ $harborUserName }}:{{ $harborPassword }}@{{ $harborDomainName }}/api/v2.0/projects'

# push docker images to {{ $.imageRepo.type }}
{{- range $_, $image := $.dockerImages }}
docker tag {{ if $image.dockerFile }}{{ $image.target }}{{ else }}{{ $image.source }}{{ end }} {{ $harborDomainName }}/{{ $image.target }}
{{- end }}

{{- range $_, $image := $.dockerImages }}
docker push {{ $harborDomainName }}/{{ $image.target }}
{{- end }}
```

## install dory services with kubernetes

```shell script
# create {{ $.dory.namespace }} namespace and pv
kubectl delete ns {{ $.dory.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
kubectl apply -f {{ $.dory.namespace }}/step01-namespace-pv.yaml

# create docker certificates
sh {{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}/docker_certs.sh
kubectl -n {{ $.dory.namespace }} create secret generic {{ $.dory.docker.dockerName }}-tls --from-file=certs/ca.crt --from-file=certs/tls.crt --from-file=certs/tls.key --dry-run=client -o yaml | kubectl apply -f -
kubectl -n {{ $.dory.namespace }} describe secret {{ $.dory.docker.dockerName }}-tls
rm -rf certs

# copy harbor certificates in docker directory
cp -rp /etc/docker/certs.d/{{ $harborDomainName }} {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}

{{- if $.dory.artifactRepo.internal.image }}
# create nexus init data, nexus init data is in a docker image
docker rm -f nexus-data-init || true
docker run -d -t --name nexus-data-init doryengine/nexus-data-init:alpine-3.15.3 cat
docker cp nexus-data-init:/nexus-data/nexus {{ $.rootDir }}/{{ $.dory.namespace }}
docker rm -f nexus-data-init
chown -R 200:200 {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
{{- end }}

# start all dory services with kubernetes
kubectl apply -f {{ $.dory.namespace }}/step02-statefulset.yaml
kubectl apply -f {{ $.dory.namespace }}/step03-service.yaml

# check dory services status
kubectl -n {{ $.dory.namespace }} get pods -o wide
```

## create project-data-alpine pod in kubernetes

```shell script
# project-data-alpine pod is used for create project directory in kuberentes
# create project-data-alpine pod in kubernetes
kubectl apply -f project-data-alpine.yaml
kubectl -n {{ $.dory.namespace }} get pods
```

## dory not config yet

2. please follow `README-kubernetes-config.md` to config dory by manual after install
