# reset docker installation

## remove all dory service when install failure

{{- if $.imageRepoInternal }}
### stop and remove {{ $.imageRepo.type }} services

```shell script
cd {{ $.rootDir }}/{{ $.imageRepo.type }}
docker-compose stop && docker-compose rm -f
```
{{- end }}

### stop and remove dory services

```shell script
cd {{ $.rootDir }}/{{ $.dory.namespace }}
docker-compose stop && docker-compose rm -f
kubectl delete namespace {{ $.dory.namespace }}
kubectl delete pv project-data-pv
```

## about dory services data

- dory services data located at: `{{ $.rootDir }}`

```shell script
# before reinstall, please remove dory services data first
rm -rf {{ $.rootDir }}/*
```
