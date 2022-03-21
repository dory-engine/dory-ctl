# doryctl changelog v0.7.7

**new features:**

- dory-core upgrade to v1.6.17
    - get kubernetes environment api resource version
    - update kubernetes project rbac permission
    - remove ingress and service if deployContainerDef removed
    - remove empty items in definition when update project definitions
    - if pipelineDef builds is empty, use default builds config
    - envK8s projectNodeSelector not required
    - reload global config dynamically
    - get kubernetes environment docker number when global config initial
    - update api api/admin/reload, support update multiple instant or dory-core global config
    - update api api/cicd/runs, support search runs by runNames
    - add api api/admin/user, support insert or update user 
    - fixed bugs: customStepDef relatedModules verify
    - fixed bugs: buildDef buildEnv verify
