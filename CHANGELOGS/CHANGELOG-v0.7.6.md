# doryctl changelog v0.7.6

**new features:**

- dory-core upgrade to v1.6.17
    - update new project kubernetes rbac permissions
        - deployment (get watch list delete scale)
        - statefulset (get watch list delete)
        - replicasset (get watch list delete scale)
        - serivce (get watch list delete)
        - pod (get watch list delete)
        - hpa (get watch list delete)
        - ingress (get watch list delete)
    - if remove deployContainerDef hpa settings, then remove hpa in kubernetes 
    - if remove deployContainerDef ingress settings, then remove ingress in kubernetes 
    - if remove deployContainerDef port settings, then remove service in kubernetes 
    - if remove deployContainerDef hpa settings, then remove hpa in kubernetes 
    - if remove componentDef port settings, then remove service in kubernetes
    - kubernetes api resource version compatibility:
        - ingress.networking.k8s.io compatible api version: v1 v1beta1
        - hpa.autoscaling compatible api version: v2 v2beta2 v2beta1 v1
        - networking.istio.io compatible api version: v1beta1 v1alpha3

