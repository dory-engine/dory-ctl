package main

import (
	"fmt"
	"github.com/alecthomas/chroma/quick"
	"os"
)

func main() {
	someSourceCode := `
version: '3'
services:
  docker:
    image: docker:20.10.7-dind
    hostname: docker
    container_name: docker
    privileged: true
    volumes:
      - ./certs:/certs
    environment:
      DOCKER_TLS_CERTDIR: ''
    command:
      - --host=tcp://0.0.0.0:2376
      - --tlsverify
      - --tlscacert=/certs/server/ca.pem
      - --tlscert=/certs/server/cert.pem
      - --tlskey=/certs/server/key.pem
    ports:
      - 2376:2376
    restart: always
`
	err := quick.Highlight(os.Stdout, someSourceCode, "yaml", "terminal", "native")
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}

	someSourceCode = `
{
    "status": "SUCCESS",
    "msg": "list projectEnv success",
    "duration": "495.347µs",
    "data": {
        "auditID": "605bfb70e4e75c57c8d12abc",
        "projectEnvs": [
            {
                "envNames": [
                    "test"
                ],
                "projectName": "test-project1"
            }
        ],
        "withAdminLog": false,
        "second": 3.2
    }
}`

	fmt.Println("############")
	err = quick.Highlight(os.Stdout, someSourceCode, "json", "terminal", "native")
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}

}
