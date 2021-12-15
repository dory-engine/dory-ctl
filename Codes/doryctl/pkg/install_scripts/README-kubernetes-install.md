# install dory with kubernetes

## summary

1. please follow `README-kubernetes-install.md` to install dory by manual
2. please follow `README-kubernetes-config.md` to config dory by manual after install
3. if install fail, please follow `README-kubernetes-reset.md` to stop all dory services and install again

## copy all scripts and config files to install root directory

```shell script
copy all scripts and config files to install root directory
mkdir -p {{ $.rootDir }}
cp -rp * {{ $.rootDir }}
```

## install {{ $.imageRepo.type }} with kubernetes

```shell script
# create {{ $.imageRepo.type }} certificates
cd {{ $.rootDir }}/{{ $.imageRepo.namespace }}
sh harbor_certs.sh
ls -alh

# install {{ $.imageRepo.type }}
cd {{ $.rootDir }}/{{ $.imageRepo.namespace }}
chmod a+x common.sh install.sh prepare
sh install.sh
ls -alh

# stop and update {{ $.imageRepo.type }} kubernetes.yml
sleep 5 && kubernetes stop && kubernetes rm -f
export HARBOR_CONFIG_ROOT_PATH=$(echo "{{ $.rootDir }}/{{ $.imageRepo.namespace }}" | sed 's#\/#\\\/#g')
sed -i "s/${HARBOR_CONFIG_ROOT_PATH}/./g" kubernetes.yml
cat kubernetes.yml

# restart {{ $.imageRepo.type }}
kubernetes up -d
sleep 10

# check {{ $.imageRepo.type }} status
kubernetes ps

# on current host and all kubernetes nodes add {{ $.imageRepo.type }} domain name in /etc/hosts
vi /etc/hosts
{{ $.hostIP }}  {{ $.imageRepo.domainName }}

# docker login to {{ $.imageRepo.type }}
docker login --username admin --password {{ $.imageRepo.password }} {{ $.imageRepo.domainName }}

# create public, hub, gcr, quay projects in {{ $.imageRepo.type }}
curl -k -X POST -H 'accept: application/json' -d '{"project_name": "public", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'accept: application/json' -d '{"project_name": "hub", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'accept: application/json' -d '{"project_name": "gcr", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'
curl -k -X POST -H 'accept: application/json' -d '{"project_name": "quay", "public": true}' 'https://admin:{{ $.imageRepo.password }}@{{ $.imageRepo.domainName }}/api/v2.0/projects'

# push docker images to {{ $.imageRepo.type }}
{{- range $_, $image := $.dockerImages }}
docker pull {{ $image.source }}
docker tag {{ $image.source }} {{ $.imageRepo.domainName }}/{{ $image.target }}
docker push {{ $.imageRepo.domainName }}/{{ $image.target }}

{{- end }}
```

## install dory services with kubernetes

```shell script
# create docker certificates
cd {{ $.rootDir }}/{{ $.dory.namespace }}/{{ $.dory.docker.dockerName }}
sh docker_certs.sh
ls -alh

# create nexus init data, nexus init data is in a docker image
cd {{ $.rootDir }}/{{ $.dory.namespace }}
docker rm -f nexus-data-init || true
docker run -d -t --name nexus-data-init dorystack/nexus-data-init:alpine-3.15.0 cat
docker cp nexus-data-init:/nexus-data/nexus .
docker rm -f nexus-data-init
chown -R 200:200 nexus
ls -alh nexus

# create dory services directory and chown
cd {{ $.rootDir }}/{{ $.dory.namespace }}
mkdir -p mongo-core-dory
chown -R 999:999 mongo-core-dory
mkdir -p dory-core/dory-data
mkdir -p dory-core/tmp
chown -R 1000:000 dory-core

# start all dory services with kubernetes
cd {{ $.rootDir }}/{{ $.dory.namespace }}
ls -alh
kubernetes stop && kubernetes rm -f && kubernetes up -d

# check dory services status
sleep 10
kubernetes ps
```

## create project-data-alpine pod in kubernetes

```shell script
# project-data-alpine pod is used for create project directory in kuberentes
# create project-data-alpine pod in kubernetes
cd {{ $.rootDir }}
kubectl apply -f project-data-alpine.yaml
kubectl -n {{ $.dory.namespace }} get pods
```
