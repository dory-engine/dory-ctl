# kubernetes环境部署要求

## 检查 /etc/timezone 和 /etc/localtime 配置

- 所有kubernetes节点必须配置 /etc/timezone 和 /etc/localtime

```shell script
# 检查 /etc/timezone 和 /etc/localtime 文件
ls -al /etc/timezone
ls -al /etc/localtime

# 更新 /etc/timezone 和 /etc/localtime 文件
timedatectl set-timezone Asia/Shanghai
echo 'Asia/Shanghai' > /etc/timezone
ls -al /etc/timezone
```

## 在kubernetes集群中创建管理token

- kubernetes管理token用于dory连接kubernetes集群并发布应用，必须在dory配置文件中设置

```shell script
# 创建管理员serviceaccount
kubectl create serviceaccount -n kube-system admin-user --dry-run=client -o yaml | kubectl apply -f -

# 创建管理员clusterrolebinding
kubectl create clusterrolebinding admin-user --clusterrole=cluster-admin --serviceaccount=kube-system:admin-user --dry-run=client -o yaml | kubectl apply -f -

# 获取kubernetes管理token
# kubernetes管理token需要在dory安装过程进行设置
kubectl -n kube-system get secret $(kubectl -n kube-system get sa admin-user -o jsonpath="{ .secrets[0].name }") -o jsonpath='{ .data.token }' | base64 -d
```

## kubernetes-dashboard

- 为了管理kubernetes中部署的应用，推荐使用`kubernetes-dashboard`
- 要了解更多，请阅读官方代码仓库README.md文档: [kubernetes-dashboard](https://github.com/kubernetes/dashboard)

- 安装:
```shell script
# 安装 kubernetes-dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.4.0/aio/deploy/recommended.yaml
```

## traefik (ingress controller)

- 要使用kubernetes的[ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)功能，必须安装ingress controller，推荐使用`traefik`
- 要了解更多，请阅读官方网站文档: [traefik](https://doc.traefik.io/traefik/)

- 在kubernetes所有master节点部署traefik: 
```shell script
# 拉取 traefik helm repo
helm repo add traefik https://helm.traefik.io/traefik
helm fetch traefik/traefik --untar

# 在kubernetes的master节点以daemonset方式部署traefik
cat << EOF > traefik.yaml
deployment:
  kind: DaemonSet
image:
  name: traefik
  tag: v2.5.4
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

# 安装traefik
kubectl create namespace traefik --dry-run=client -o yaml | kubectl apply -f -
helm install -n traefik traefik traefik/ -f traefik.yaml

# 检查安装情况
helm -n traefik list
kubectl -n traefik get pods -o wide
kubectl -n traefik get services -o wide
```

## metrics-server

- 为了使用kubernetes的水平扩展缩容功能[horizontal pod autoscale](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)，必须安装`metrics-server`
- 要了解更多，请阅读官方代码仓库README.md文档: [metrics-server](https://github.com/kubernetes-sigs/metrics-server)

- install:
```shell script
# 拉取镜像
docker pull registry.aliyuncs.com/google_containers/metrics-server:v0.5.2
docker tag registry.aliyuncs.com/google_containers/metrics-server:v0.5.2 k8s.gcr.io/metrics-server/metrics-server:v0.5.2

# 获取metrics-server安装yaml
curl -O -L https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.5.2/components.yaml
# 添加--kubelet-insecure-tls参数
sed -i 's/- args:/- args:\n        - --kubelet-insecure-tls/g' components.yaml
# 安装metrics-server
kubectl apply -f components.yaml

# 等待metrics-server正常
kubectl -n kube-system get pods -l=k8s-app=metrics-server

# 查看pod的metrics
kubectl top pods -A

```
