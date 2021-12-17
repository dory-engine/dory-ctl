# kubernetes prerequisite

## kubernetes-dashboard

- to manager your project pods in kubernetes, we recommend to use `kubernetes-dashboard`
- read official repository README.md to learn more: [kubernetes-dashboard](https://github.com/kubernetes/dashboard)

- install:
```shell script
# install kubernetes-dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.4.0/aio/deploy/recommended.yaml
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
# install metrics-server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```
