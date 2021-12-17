# install dory with kubernetes

## summary

1. please follow `README-kubernetes-install.md` to install dory by manual
2. please follow `README-kubernetes-config.md` to config dory by manual after install
3. if install fail, please follow `README-kubernetes-reset.md` to stop all dory services and install again

## create install root directories

```shell script
# create {{ $.imageRepo.type }} root directory
mkdir -p {{ $.rootDir }}/{{ $.imageRepo.namespace }}/database
mkdir -p {{ $.rootDir }}/{{ $.imageRepo.namespace }}/jobservice
mkdir -p {{ $.rootDir }}/{{ $.imageRepo.namespace }}/redis
mkdir -p {{ $.rootDir }}/{{ $.imageRepo.namespace }}/registry
chown -R 999:999 {{ $.rootDir }}/{{ $.imageRepo.namespace }}/database
chown -R 10000:10000 {{ $.rootDir }}/{{ $.imageRepo.namespace }}/jobservice
chown -R 999:999 {{ $.rootDir }}/{{ $.imageRepo.namespace }}/redis
chown -R 10000:10000 {{ $.rootDir }}/{{ $.imageRepo.namespace }}/registry
ls -alh {{ $.rootDir }}/{{ $.imageRepo.namespace }}

# create dory root directory
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/dory-data
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/dory-core/tmp
mkdir -p {{ $.rootDir }}/{{ $.imageRepo.namespace }}
cp -rp {{ $.dory.namespace }}/{{ $.dory.docker.dockerName }} {{ $.rootDir }}/{{ $.dory.namespace }}/
cp -rp {{ $.dory.namespace }}/dory-core {{ $.rootDir }}/{{ $.dory.namespace }}/
chown -R 1000:1000 {{ $.rootDir }}/{{ $.dory.namespace }}/dory-core
mkdir -p {{ $.rootDir }}/{{ $.dory.namespace }}/mongo-core-dory
chown -R 999:999 {{ $.rootDir }}/{{ $.dory.namespace }}/mongo-core-dory
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}

```

## install {{ $.imageRepo.type }} with kubernetes

```shell script
# create {{ $.imageRepo.type }} namespace and pv
kubectl delete ns {{ $.imageRepo.namespace }}
kubectl delete pv {{ $.imageRepo.namespace }}-pv
kubectl apply -f {{ $.imageRepo.namespace }}/step01-namespace-pv.yaml

# install {{ $.imageRepo.type }}
helm install -n {{ $.imageRepo.namespace }} {{ $.imageRepo.namespace }} {{ $.imageRepo.type }}
helm -n {{ $.imageRepo.namespace }} list

# waiting for all {{ $.imageRepo.type }} services ready
kubectl -n {{ $.imageRepo.namespace }} get pods -o wide

# get {{ $.imageRepo.type }} and copy to /etc/docker/certs.d
sh {{ $.imageRepo.namespace }}/harbor_update_docker_certs.sh
ls -alh /etc/docker/certs.d/{{ $.imageRepo.domainName }}

# on current host and all kubernetes nodes add {{ $.imageRepo.type }} domain name in /etc/hosts
vi /etc/hosts
{{ $.hostIP }}  {{ $.imageRepo.domainName }}

# docker login to {{ $.imageRepo.type }}
docker login --username admin --password {{ $.imageRepo.password }} {{ $.imageRepo.domainName }}

# create public, hub, gcr, quay projects in {{ $.imageRepo.type }}
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "public", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "hub", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "gcr", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'Content-Type: application/json' -d '{"project_name": "quay", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'

# push docker images to {{ $.imageRepo.type }}
{{- range $_, $image := $.dockerImages }}
docker tag {{ $image.source }} {{ $.imageRepo.domainName }}/{{ $image.target }}
{{- end }}

{{- range $_, $image := $.dockerImages }}
docker push {{ $.imageRepo.domainName }}/{{ $image.target }}
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
cp -rp /etc/docker/certs.d/{{ $.imageRepo.domainName }} {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}


# create nexus init data, nexus init data is in a docker image
docker rm -f nexus-data-init || true
docker run -d -t --name nexus-data-init doryengine/nexus-data-init:alpine-3.15.0 cat
docker cp nexus-data-init:/nexus-data/nexus {{ $.rootDir }}/{{ $.dory.namespace }}
docker rm -f nexus-data-init
chown -R 200:200 {{ $.rootDir }}/{{ $.dory.namespace }}/nexus
ls -alh {{ $.rootDir }}/{{ $.dory.namespace }}/nexus

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
