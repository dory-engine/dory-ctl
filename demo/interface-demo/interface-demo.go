package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
)

func removeItem(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		vv := reflect.ValueOf(v)
		switch v.(type) {
		case int, float32, float64:
			if vv.IsZero() {
				delete(m, k)
			}
		case bool:
			if vv.Bool() == false {
				delete(m, k)
			}
		case string:
			if vv.String() == "" {
				delete(m, k)
			}
		}
		if !vv.IsValid() {
			delete(m, k)
		}
		if vv.Kind() == reflect.Slice {
			if vv.Len() == 0 {
				delete(m, k)
			} else {
				var isMap bool
				var x []map[string]interface{}
				for i := 0; i < vv.Len(); i++ {
					vvv := reflect.ValueOf(vv.Index(i))
					if vvv.Kind() == reflect.Map {
						vm, ok := vv.Index(i).Interface().(map[string]interface{})
						if ok {
							isMap = true
							v3 := removeItem(vm)
							x = append(x, v3)
						}
					} else if vvv.Kind() == reflect.Struct {
						vm, ok := vv.Index(i).Interface().(map[string]interface{})
						if ok {
							isMap = true
							v3 := removeItem(vm)
							x = append(x, v3)
						}
					}
				}
				if isMap {
					m[k] = x
				}
			}
		}
		if vv.Kind() == reflect.Struct {
			v2 := removeItem(v.(map[string]interface{}))
			if len(v2) == 0 {
				delete(m, k)
			} else {
				m[k] = v2
			}
		} else if vv.Kind() == reflect.Map {
			v2 := removeItem(v.(map[string]interface{}))
			if len(v2) == 0 {
				delete(m, k)
			} else {
				m[k] = v2
			}
		}
	}
	m2 := m
	return m2
}

func main() {
	m := map[string]interface{}{
		"string_foo":   "foo",
		"string_empty": "",
		"int_zero":     0,
		"int":          1,
		"bool_false":   false,
		"bool_true":    true,
		"nil":          nil,
		"array_empty":  []string{"a"},
		"array_maps": []map[string]interface{}{
			{
				"string_foo":   "foo",
				"string_empty": "",
				"int_zero":     0,
				"int":          1,
				"bool_false":   false,
				"bool_true":    true,
				"nil":          nil,
				"array":        []string{},
			},
			{
				"string_foo":   "foo",
				"string_empty": "",
				"int_zero":     0,
				"int":          1,
				"bool_false":   false,
				"bool_true":    true,
				"nil":          nil,
				"array":        []string{},
			},
		},
		"map": map[string]interface{}{
			"string_foo":   "foo",
			"string_empty": "",
			"int_zero":     0,
			"int":          1,
			"bool_false":   false,
			"bool_true":    true,
			"nil":          nil,
			"array_empty":  []string{},
			"array_maps": []map[string]interface{}{
				{
					"string_foo":   "foo",
					"string_empty": "",
					"int_zero":     0,
					"int":          1,
					"bool_false":   false,
					"bool_true":    true,
					"nil":          nil,
					"array":        []string{},
				},
				{
					"string_foo":   "foo",
					"string_empty": "",
					"int_zero":     0,
					"int":          1,
					"bool_false":   false,
					"bool_true":    true,
					"nil":          nil,
					"array":        []string{},
				},
			},
		},
	}
	m2 := removeItem(m)
	bs, _ := yaml.Marshal(m2)
	fmt.Println(string(bs))
	fmt.Println("##############")

	strYaml := `
deploys:
- deployName: tp1-gin-demo
  relatedPackage: tp1-gin-demo
  deployImageTag: ""
  deployLabels: {}
  deploySessionAffinityTimeoutSeconds: 0
  deployNodePorts:
    - port: 8000
      nodePort: 30103
      protocol: http
  deployLocalPorts: []
  deployReplicas: 1
  hpaConfig:
    maxReplicas: 2
    memoryAverageValue: 100Mi
    memoryAverageRequestPercent: 0
    cpuAverageValue: 100m
    cpuAverageRequestPercent: 0
  deployEnvs: []
  deployCommand: sh -c "./tp1-gin-demo 2>&1 | sed \"s/^/[$(hostname)] /\" | tee -a /tp1-gin-demo/logs/tp1-gin-demo.logs"
  deployCmd: []
  deployResources:
    memoryRequest: 10Mi
    memoryLimit: 100Mi
    cpuRequest: "0.02"
    cpuLimit: "0.1"
  deployVolumes:
    - pathInPod: /tp1-gin-demo/logs
      pathInPv: tp1-gin-demo/logs
      pvc: ""
  deployHealthCheck:
    checkPort: 0
    httpGet:
        path: /
        port: 8000
        httpHeaders: []
    readinessDelaySeconds: 15
    readinessPeriodSeconds: 5
    livenessDelaySeconds: 150
    livenessPeriodSeconds: 30
  dependServices: []
  hostAliases: []
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  deployConfigSettings:
    - Codes/Backend/tp1-gin-demo/config1/:tp1-gin-demo/config1/
    - Codes/Backend/tp1-gin-demo/config2/:tp1-gin-demo/config2/
- deployName: tp1-go-demo
  relatedPackage: tp1-go-demo
  deployImageTag: ""
  deployLabels: {}
  deploySessionAffinityTimeoutSeconds: 0
  deployNodePorts:
    - port: 8000
      nodePort: 30102
      protocol: http
  deployLocalPorts: []
  deployReplicas: 1
  hpaConfig:
    maxReplicas: 0
    memoryAverageValue: ""
    memoryAverageRequestPercent: 0
    cpuAverageValue: ""
    cpuAverageRequestPercent: 0
  deployEnvs: []
  deployCommand: sh -c "./tp1-go-demo 2>&1 | sed \"s/^/[$(hostname)] /\" | tee -a /tp1-go-demo/logs/tp1-go-demo.logs"
  deployCmd: []
  deployResources:
    memoryRequest: 10Mi
    memoryLimit: 150Mi
    cpuRequest: "0.02"
    cpuLimit: "0.1"
  deployVolumes:
    - pathInPod: /tp1-go-demo/logs
      pathInPv: tp1-go-demo/logs
      pvc: ""
  deployHealthCheck:
    checkPort: 0
    httpGet:
        path: /
        port: 8000
        httpHeaders: []
    readinessDelaySeconds: 15
    readinessPeriodSeconds: 5
    livenessDelaySeconds: 150
    livenessPeriodSeconds: 30
  dependServices:
    - dependName: tp1-node-demo
      dependPort: 3000
      dependType: TCP
    - dependName: tp1-python-demo
      dependPort: 3000
      dependType: TCP
  hostAliases: []
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  deployConfigSettings: []
- deployName: tp1-node-demo
  relatedPackage: tp1-node-demo
  deployImageTag: ""
  deployLabels: {}
  deploySessionAffinityTimeoutSeconds: 0
  deployNodePorts: []
  deployLocalPorts:
    - port: 3000
      protocol: http
      ingress:
        domainName: ""
        pathPrefix: ""
  deployReplicas: 1
  hpaConfig:
    maxReplicas: 0
    memoryAverageValue: ""
    memoryAverageRequestPercent: 0
    cpuAverageValue: ""
    cpuAverageRequestPercent: 0
  deployEnvs: []
  deployCommand: sh -c "node index.js 2>&1 | sed \"s/^/[$(hostname)] /\" | tee -a /tp1-node-demo/logs/tp1-node-demo.logs"
  deployCmd: []
  deployResources:
    memoryRequest: 10Mi
    memoryLimit: 150Mi
    cpuRequest: "0.02"
    cpuLimit: "0.1"
  deployVolumes:
    - pathInPod: /tp1-node-demo/logs
      pathInPv: tp1-node-demo/logs
      pvc: ""
  deployHealthCheck:
    checkPort: 0
    httpGet:
        path: /
        port: 3000
        httpHeaders: []
    readinessDelaySeconds: 15
    readinessPeriodSeconds: 5
    livenessDelaySeconds: 150
    livenessPeriodSeconds: 30
  dependServices: []
  hostAliases: []
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  deployConfigSettings: []
- deployName: tp1-python-demo
  relatedPackage: tp1-python-demo
  deployImageTag: ""
  deployLabels: {}
  deploySessionAffinityTimeoutSeconds: 0
  deployNodePorts: []
  deployLocalPorts:
    - port: 3000
      protocol: http
      ingress:
        domainName: ""
        pathPrefix: ""
  deployReplicas: 1
  hpaConfig:
    maxReplicas: 0
    memoryAverageValue: ""
    memoryAverageRequestPercent: 0
    cpuAverageValue: ""
    cpuAverageRequestPercent: 0
  deployEnvs: []
  deployCommand: sh -c "python3 main.py 2>&1 | sed \"s/^/[$(hostname)] /\" | tee -a /tp1-python-demo/logs/tp1-python-demo.logs"
  deployCmd: []
  deployResources:
    memoryRequest: 10Mi
    memoryLimit: 150Mi
    cpuRequest: "0.02"
    cpuLimit: "0.1"
  deployVolumes:
    - pathInPod: /tp1-python-demo/logs
      pathInPv: tp1-python-demo/logs
      pvc: ""
  deployHealthCheck:
    checkPort: 0
    httpGet:
        path: /
        port: 3000
        httpHeaders: []
    readinessDelaySeconds: 15
    readinessPeriodSeconds: 5
    livenessDelaySeconds: 150
    livenessPeriodSeconds: 30
  dependServices: []
  hostAliases: []
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  deployConfigSettings: []
- deployName: tp1-spring-demo
  relatedPackage: tp1-spring-demo
  deployImageTag: ""
  deployLabels: {}
  deploySessionAffinityTimeoutSeconds: 0
  deployNodePorts: []
  deployLocalPorts:
    - port: 9000
      protocol: http
      ingress:
        domainName: demo.test-project1.local
        pathPrefix: /spring/
  deployReplicas: 1
  hpaConfig:
    maxReplicas: 0
    memoryAverageValue: ""
    memoryAverageRequestPercent: 0
    cpuAverageValue: ""
    cpuAverageRequestPercent: 0
  deployEnvs:
    - JAVA_OPTS=-Xms256m -Xmx256m
  deployCommand: sh -c "java -jar example.smallest-0.0.1-SNAPSHOT.war 2>&1 | sed \"s/^/[$(hostname)] /\" | tee -a /tp1-spring-demo/logs/tp1-spring-demo.logs"
  deployCmd: []
  deployResources:
    memoryRequest: 10Mi
    memoryLimit: 250Mi
    cpuRequest: "0.05"
    cpuLimit: "0.25"
  deployVolumes:
    - pathInPod: /tp1-spring-demo/logs
      pathInPv: tp1-spring-demo/logs
      pvc: ""
  deployHealthCheck:
    checkPort: 0
    httpGet:
        path: /
        port: 9000
        httpHeaders: []
    readinessDelaySeconds: 15
    readinessPeriodSeconds: 5
    livenessDelaySeconds: 150
    livenessPeriodSeconds: 30
  dependServices: []
  hostAliases: []
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  deployConfigSettings: []
`
	var mapYaml map[string]interface{}
	err := yaml.Unmarshal([]byte(strYaml), &mapYaml)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		return
	}
	mapYaml = removeItem(mapYaml)
	bs, _ = yaml.Marshal(mapYaml)
	fmt.Println(string(bs))

}
