# reset kubernetes installation

## remove all dory service when install failure

{{- if $.imageRepoInternal }}
### stop and remove {{ $.imageRepo.type }} services

```shell script
helm -n {{ $.imageRepo.internal.namespace }} uninstall {{ $.imageRepo.internal.namespace }}
```
{{- end }}

### stop and remove dory services

```shell script
kubectl delete namespace {{ $.dory.namespace }}
kubectl delete pv {{ $.dory.namespace }}-pv
{{- if $.imageRepoInternal }}
kubectl delete namespace {{ $.imageRepo.internal.namespace }}
kubectl delete pv {{ $.imageRepo.internal.namespace }}-pv
{{- end }}
kubectl delete pv project-data-pv
```

## about dory services data

- dory services data located at: `{{ $.rootDir }}`

```shell script
# before reinstall, please remove dory services data first
rm -rf {{ $.rootDir }}/*
```
