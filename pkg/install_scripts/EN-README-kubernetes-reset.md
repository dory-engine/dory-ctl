# reset kubernetes installation

## remove all dory service when install failure

{{- if $.imageRepo.internal.domainName }}
### stop and remove {{ $.imageRepo.type }} services

```shell script
helm -n {{ $.imageRepo.internal.namespace }} uninstall {{ $.imageRepo.internal.namespace }}
```
{{- end }}

### stop and remove dory services

```shell script
# cd to readme directory
kubectl delete namespace {{ $.dory.namespace }}
kubectl delete namespace {{ $.imageRepo.internal.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
kubectl delete pv {{ $.imageRepo.internal.namespace }}-pv
kubectl delete pv project-data-pv
```

## about dory services data

- dory services data located at: `{{ $.rootDir }}`

```shell script
# before reinstall, please remove dory services data first
rm -rf {{ $.rootDir }}/*
```
