# doryctl is a command line toolkit to install and control `Dory-Engine`

![](docs/images/dory-icon.png)

The official website: [https://doryengine.com](https://doryengine.com)

## What is `Dory-Engine`

- `Dory-Engine` is an engine to make your application to Cloud-Native infrastructure extremely easy. 

## Install Dory-Engine by doryctl

- Now we can use doryctl to install `Dory-Engine` with `docker-compose` (for test usage) or `kubernetes`(for production usage, recommended)

```shell script
##############################
please follow these steps to install dory-core with docker:

# 1. check prerequisite for install with docker
doryctl install check --mode docker

# 2. pull relative docker images from docker hub
doryctl install pull

# 3. print docker install mode config settings
doryctl install print --mode docker > install-config-docker.yaml

# 4. update install config file by manual
vi install-config-docker.yaml

# 5. install dory with docker
doryctl install run -o readme-install-docker -f install-config-docker.yaml

##############################
# please follow these steps to install dory-core with kubernetes:

# 1. check prerequisite for install with kubernetes
doryctl install check --mode kubernetes

# 2. pull relative docker images from docker hub
doryctl install pull

# 3. print kubernetes install mode config settings
doryctl install print --mode kubernetes > install-config-kubernetes.yaml

# 4. update install config file by manual
vi install-config-kubernetes.yaml

# 5. install dory with kubernetes
doryctl install run -o readme-install-kubernetes -f install-config-kubernetes.yaml
```

- For more detail:

```shell script
doryctl -h
```
