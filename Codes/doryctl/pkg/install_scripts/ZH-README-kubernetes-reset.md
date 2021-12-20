# 清除kubernetes方式部署的服务

## 当安装出现异常的情况下，清除所有dory服务

### 停止并清除 {{ $.imageRepo.type }} 服务

```shell script
helm -n {{ $.imageRepo.namespace }} uninstall {{ $.imageRepo.namespace }}
```

### 停止并清除所有 dory 服务

```shell script
# cd to readme directory
kubectl delete namespace {{ $.dory.namespace }}
kubectl delete namespace {{ $.imageRepo.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
kubectl delete pv {{ $.imageRepo.namespace }}-pv
kubectl delete pv project-data-pv
```

## 所有dory组件的数据存放位置

- 所有dory组件的数据存放在: `{{ $.rootDir }}`

```shell script
# 重新安装前，请清理dory组件数据
rm -rf {{ $.rootDir }}/*
```
