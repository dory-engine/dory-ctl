# doryctl changelog v0.8.1

**new features:**

- dory-core upgrade to v1.7.0
    - fixed bugs: step retry and timeout abort return error
    - stage add isContainer param, enhance performance of pipeline runtime
    - step timeout count begin until docker image pulled
    - support abort run while create container and pulling image
    - remove packr, use golang embed filesystem
    - remove config.yaml dockerHubUrl settings
    - spring-demo upgrade to jdk11
    - node-demo and vue-demo upgrade to npm-node15
    - maven build environments update to: jdk8 jdk11 jdk17
    - gradle build environments update to: jdk8 jdk11 jdk17
    - update kubernetes namespace harobr and nexus secrets template
    - fixed bugs: kubernetes installation way run stuck when npm install, remove vue-demo package-lock.json
    - config.yaml add git repo callback url

- doryctl upgrade to v0.8.1
    - support use external git repository, artifact repository, image repository
    - upgrade build environments docker images
    - upgrade dory services docker images
    - ignore harbor project create error if deploy harbor external way
    - update readme files
    - add external git repository settings: dory URL callback by git repository webhook
