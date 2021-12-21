package main

import (
	"encoding/json"
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
				x := []map[string]interface{}{}
				for i := 0; i < vv.Len(); i++ {
					vvv := reflect.ValueOf(vv.Index(i))
					fmt.Println(vvv)
					if vvv.Kind() == reflect.String {
					} else if vvv.Kind() == reflect.Int {
					} else if vvv.Kind() == reflect.Float32 {
					} else if vvv.Kind() == reflect.Float64 {
					} else if vvv.Kind() == reflect.Bool {
					} else if vvv.Kind() == reflect.Struct {
						isMap = true
						v3 := removeItem(vv.Index(i).Interface().(map[string]interface{}))
						x = append(x, v3)
					} else if vvv.Kind() == reflect.Map {
						isMap = true
						v3 := removeItem(vv.Index(i).Interface().(map[string]interface{}))
						x = append(x, v3)
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
	data, _ := json.Marshal(m2)
	fmt.Println(string(data))

	strYaml := `
deployName: tp1-spring-demo
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
	bs, _ := yaml.Marshal(mapYaml)
	fmt.Println(string(bs))

}
