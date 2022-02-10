# doryctl changelog v0.7.0

**new features:**

- dory-core upgrade to v1.6.9
    - customStepConf support set step status as failed in dory-param-input file
    - customStepConf support save output files in assets directory
    - customStepConf support save output files in tar file
    - customStepConf support when dockerVolumes is empty, mount repository in docker executor by default
    - add api: /public/about
    
- dory-dashboard upgrade to v1.6.3
    - customStep add tarFile and outputFiles support
    - add api support: /public/about
    - customStepConf dockerVolumes and dockerEnvs webUI upgraded

- doryctl upgrade
    - dory-core config.yaml add customStepOutputDir, customStepOutputUri, dockerParamInputFileName, dockerParamOutputFileName, dockerOutputFileDir options
    - add login command
    - add logout command
    - add project get command
    - add pipeline get command
    - doryctl log with verbose mode
    
**fixed bug**

- dory-core upgrade to v1.6.9
    - can't delete customStepConf bug fixed 
