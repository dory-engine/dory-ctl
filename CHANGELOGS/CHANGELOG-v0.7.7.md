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

- doryctl upgrade to v0.7.7
    - update kubernetes install readme
    - print yaml indent
    - now doryctl support manage project definition resource with def command
    - add command: def get (get project definitions)
    - add command: def apply (apply project definitions)
    - add command: def delete (delete modules from project definitions)
    - add command: def patch (patch project definitions)
    - add command: def clone (clone project definitions modules to another environments)
    - now doryctl support manage server configurations with admin command
    - add command: admin get (get configurations, admin permission required)
    - add command: admin apply (apply configurations, admin permission required)
    - add command: admin delete (delete configurations, admin permission required)
    - now doryctl all command support completion
    - doryctl version command support show server version
