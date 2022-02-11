# doryctl changelog v0.7.0

**new features:**

- dory-core upgrade to v1.6.10
    - customStepConf support set step status as failed in dory-param-input file
    - customStepConf support save output files in assets directory
    - customStepConf support save output files in tar file
    - customStepConf support when dockerVolumes is empty, mount repository in docker executor by default
    - add api: /public/about
    - check token in request header or query automatically
    - name id string length unlimited
    
- dory-dashboard upgrade to v1.6.3
    - customStep add tarFile and outputFiles support
    - add api support: /public/about
    - customStepConf dockerVolumes and dockerEnvs webUI upgraded

- doryctl upgrade
    - dory-core config.yaml file add customStepOutputDir, customStepOutputUri, dockerParamInputFileName, dockerParamOutputFileName, dockerOutputFileDir options
    - log with verbose mode
    - add command: login
    - add command: logout
    - add command: project get
    - add command: pipeline get
    - add command: pipeline execute
    - add command: run get
    - add command: run logs

**fixed bug**

- dory-core upgrade to v1.6.10
    - can't delete customStepConf bug fixed 
    
***bugs***
- doryctl run logs command: finished run read message stuck until timeout 

