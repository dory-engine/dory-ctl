# kubernetes prerequisite

## caution about docker install mode 

- docker mode install machine and kubernetes master can't be the same machine, because in docker mode harbor services and kubernetes master ingress controller use the same 443 TLS port, port conflicts will cause dory connect harbor failed
- make sure you have 2 hosts for docker mode install 

## hardware requirement

- cpus: 2 cores
- memory: 8G
- storage: 40G

## check /etc/timezone and /etc/localtime

- all kubernetes nodes should have /etc/timezone and /etc/localtime files

```shell script
# check /etc/timezone and /etc/localtime files
ls -al /etc/timezone
ls -al /etc/localtime

# update /etc/timezone and /etc/localtime files
timedatectl set-timezone Asia/Shanghai
echo 'Asia/Shanghai' > /etc/timezone
ls -al /etc/timezone
```

## create kubernetes admin token

- kubernetes admin token is for dory to deploy project applications in kubernetes cluster, you must set it in dory's config file

```shell script
# create kubernetes admin serviceaccount
kubectl create serviceaccount -n kube-system admin-user --dry-run=client -o yaml | kubectl apply -f -

# create kubernetes admin clusterrolebinding
kubectl create clusterrolebinding admin-user --clusterrole=cluster-admin --serviceaccount=kube-system:admin-user --dry-run=client -o yaml | kubectl apply -f -

# get kubernetes admin token
# kubernetes token is for dory installation config
kubectl -n kube-system get secret $(kubectl -n kube-system get sa admin-user -o jsonpath="{ .secrets[0].name }") -o jsonpath='{ .data.token }' | base64 -d
```

## kubernetes-dashboard

- to manager your project pods in kubernetes, we recommend to use `kubernetes-dashboard`
- read official repository README.md to learn more: [kubernetes-dashboard](https://github.com/kubernetes/dashboard)

- install:
```shell script
# install kubernetes-dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.5.1/aio/deploy/recommended.yaml
```

## traefik (ingress controller)

- to use kubernetes [ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/), you need to install an ingress controller, we recommend to use `traefik`
- read official website to learn more: [traefik](https://doc.traefik.io/traefik/)

- install traefik in kubernetes master nodes:
```shell script
# fetch traefik helm repo
helm repo add traefik https://helm.traefik.io/traefik
helm fetch traefik/traefik --untar

# install traefik in kubernetes as daemonset on master nodes
cat << EOF > traefik.yaml
deployment:
  kind: DaemonSet
image:
  name: traefik
  tag: v2.6.3
pilot:
  enabled: true
experimental:
  plugins:
    enabled: true
ports:
  web:
    hostPort: 80
  websecure:
    hostPort: 443
service:
  type: ClusterIP
nodeSelector:
  node-role.kubernetes.io/master: ""
EOF

# install traefik
kubectl create namespace traefik --dry-run=client -o yaml | kubectl apply -f -
helm install -n traefik traefik traefik/ -f traefik.yaml

# check install
helm -n traefik list
kubectl -n traefik get pods -o wide
kubectl -n traefik get services -o wide
```

## metrics-server

- to use kubernetes [horizontal pod autoscale](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/), must install `metrics-server`
- read official repository README.md to learn more: [metrics-server](https://github.com/kubernetes-sigs/metrics-server)

- install:
```shell script
# pull docker image
docker pull registry.aliyuncs.com/google_containers/metrics-server:v0.6.1
docker tag registry.aliyuncs.com/google_containers/metrics-server:v0.6.1 k8s.gcr.io/metrics-server/metrics-server:v0.6.1

# get metrics-server install yaml
curl -O -L https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.6.1/components.yaml
# add --kubelet-insecure-tls args
sed -i 's/- args:/- args:\n        - --kubelet-insecure-tls/g' components.yaml
# install metrics-server
kubectl apply -f components.yaml


# waiting for metrics-server ready
kubectl -n kube-system get pods -l=k8s-app=metrics-server

# get pods metrics
kubectl top pods -A
```
