# reset kubernetes installation

## remove all dory service when install failure

### stop and remove {{ $.imageRepo.type }} services

```shell script
helm -n {{ $.imageRepo.namespace }} uninstall {{ $.imageRepo.namespace }}
```

### stop and remove dory services

```shell script
# cd to readme directory
kubectl delete namespace {{ $.dory.namespace }}
kubectl delete namespace {{ $.imageRepo.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
kubectl delete pv {{ $.imageRepo.namespace }}-pv
kubectl delete pv project-data-pv
```

## about dory services data

- dory services data located at: `{{ $.rootDir }}`

```shell script
# before reinstall, please remove dory services data first
rm -rf {{ $.rootDir }}/*
```
