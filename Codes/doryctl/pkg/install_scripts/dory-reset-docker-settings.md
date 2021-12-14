# reset docker installation

## remove all dory service when install failure

### stop and remove {{ $.imageRepo.type }} services

```shell script
cd {{ $.rootDir }}/{{ $.imageRepo.type }}
docker-compose stop && docker-compose rm -f
```

### stop and remove dory services

```shell script
cd {{ $.rootDir }}/{{ $.dory.namespace }}
docker-compose stop && docker-compose rm -f
```

## about dory services data

- dory services data located at: `{{ $.rootDir }}`
