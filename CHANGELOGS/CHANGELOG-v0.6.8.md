# doryctl changelog v0.6.8

**new features:**

- dory-core upgrade to v1.6.8
    - customStepConf support set step status as failed in dory-param-input file
    - customStepConf support save output files in assets directory
    - customStepConf support save output files in tar file
    - customStepConf support when dockerVolumes is empty, mount repository in docker executor by default
    - add api: /public/about
    
- doryctl upgrade
    - config.yaml add customStepOutputDir, customStepOutputUri, dockerParamInputFileName, dockerParamOutputFileName, dockerOutputFileDir options
