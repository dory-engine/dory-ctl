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
deployName: tp1-spring-demo
port: 9000
protocol: http
httpSettings:
  matchHeaders: []
  gateway:
    rewriteUri: /spring
    matchUris: []
    matchDefault: false
  timeout: ""
  retries:
    retryOn: ""
    attempts: 0
    perTryTimeout: ""
  mirror:
    host: ""
    port: 0
    subset: ""
    mirrorPercent: 0
  corsPolicy:
    allowOrigins:
      - gateway:
          rewriteUri: /spring
          matchUris:
            - gateway:
                rewriteUri: /spring
                matchUris:
                - rewriteUri: /spring
                  matchUris:
                  - rewriteUri: /spring
                matchDefault: false
          matchDefault: false
    allowMethods: []
    allowCredentials: false
    allowHeaders: []
    exposeHeaders:
      - int: 0
        x:
          - int: 0
    maxAge: ""
  trafficPolicyEnable: false
  loadBalancer:
    loadBalancerEnable: false
    simple: ""
    consistentHash:
      consistentHashEnable: false
      httpHeaderName: ""
      httpCookie:
        name: ""
        path: ""
        ttl: ""
      useSourceIp: false
      httpQueryParameterName: ""
  connectionPool:
    connectionPoolEnable: false
    tcp:
      tcpEnable: false
      maxConnections: 0
      connectTimeout: ""
    http:
      httpEnable: false
      http1MaxPendingRequests: 0
      http2MaxRequests: 0
      maxRequestsPerConnection: 0
      maxRetries: 0
      idleTimeout: ""
  outlierDetection:
    outlierDetectionEnable: false
    consecutiveGatewayErrors: 1.01
    consecutive5xxErrors: 0.1
    interval: ""
    baseEjectionTime: ""
    maxEjectionPercent: 0
    minHealthPercent: 0
tcpSettings:
  sourceServiceNames: [""]
labelName: ""
localLabelConfig:
  labelDefault: ""
  labelNew: ""
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
