package pkg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (ic *InstallConfig) VerifyInstallConfig() error {
	var err error
	errInfo := fmt.Sprintf("verify install config error")

	var fieldName, fieldValue string

	fieldName = "installMode"
	fieldValue = ic.InstallMode
	if fieldValue != "docker" && fieldValue != "kubernetes" {
		err = fmt.Errorf("%s: %s %s format error: must be docker or kubernetes", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "rootDir"
	fieldValue = ic.RootDir
	if !strings.HasPrefix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
		return err
	}
	if strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "dory.namespace"
	fieldValue = ic.Dory.Namespace
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "dory.gitRepo.type"
	fieldValue = ic.Dory.GitRepo.Type
	if ic.Dory.GitRepo.Type != "gitea" && ic.Dory.GitRepo.Type != "gitlab" {
		err = fmt.Errorf("%s: %s %s format error: must be gitea or gitlab", errInfo, fieldName, fieldValue)
		return err
	}

	if ic.Dory.GitRepo.Internal.Image == "" && ic.Dory.GitRepo.External.ViewUrl == "" {
		err = fmt.Errorf("%s: dory.gitRepo.internal and dory.gitRepo.external both empty", errInfo)
		return err
	}

	if ic.Dory.GitRepo.Internal.Image != "" && ic.Dory.GitRepo.External.ViewUrl != "" {
		err = fmt.Errorf("%s: dory.gitRepo.internal and dory.gitRepo.external can not set at the same time", errInfo)
		return err
	}

	if ic.Dory.GitRepo.Internal.Image != "" {
		fieldName = "dory.gitRepo.internal.imageDB"
		fieldValue = ic.Dory.GitRepo.Internal.ImageDB
		if ic.Dory.GitRepo.Type == "gitea" && ic.Dory.GitRepo.Internal.ImageDB == "" {
			err = fmt.Errorf("%s: %s %s format error: gitea imageDB can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
	}

	fieldName = "dory.artifactRepo.type"
	fieldValue = ic.Dory.ArtifactRepo.Type
	if ic.Dory.ArtifactRepo.Type != "nexus" {
		err = fmt.Errorf("%s: %s %s format error: must be nexus", errInfo, fieldName, fieldValue)
		return err
	}

	if ic.Dory.ArtifactRepo.Internal.Image == "" && ic.Dory.ArtifactRepo.External.ViewUrl == "" {
		err = fmt.Errorf("%s: dory.artifactRepo.internal and dory.artifactRepo.external both empty", errInfo)
		return err
	}

	if ic.Dory.ArtifactRepo.Internal.Image != "" && ic.Dory.ArtifactRepo.External.ViewUrl != "" {
		err = fmt.Errorf("%s: dory.artifactRepo.internal and dory.artifactRepo.external can not set at the same time", errInfo)
		return err
	}

	fieldName = "imageRepo.type"
	fieldValue = ic.ImageRepo.Type
	if ic.ImageRepo.Type != "harbor" {
		err = fmt.Errorf("%s: %s %s format error: must be harbor", errInfo, fieldName, fieldValue)
		return err
	}

	if ic.ImageRepo.Internal.DomainName == "" && ic.ImageRepo.External.ViewUrl == "" {
		err = fmt.Errorf("%s: imageRepo.internal and imageRepo.external both empty", errInfo)
		return err
	}

	if ic.ImageRepo.Internal.DomainName != "" && ic.ImageRepo.External.ViewUrl != "" {
		err = fmt.Errorf("%s: imageRepo.internal and imageRepo.external can not set at the same time", errInfo)
		return err
	}

	if ic.ImageRepo.Internal.DomainName != "" {
		fieldName = "imageRepo.internal.namespace"
		fieldValue = ic.ImageRepo.Internal.Namespace
		if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
			err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
			return err
		}

		if ic.InstallMode == "docker" {
			fieldName = "imageRepo.internal.certsDir"
			fieldValue = ic.ImageRepo.Internal.CertsDir
			if fieldValue == "" {
				err = fmt.Errorf("%s: %s %s format error: installMode is docker, imageRepo.internal.certsDir can not be empty", errInfo, fieldName, fieldValue)
				return err
			}
			if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
				err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
				return err
			}

			fieldName = "imageRepo.internal.dataDir"
			fieldValue = ic.ImageRepo.Internal.DataDir
			if fieldValue == "" {
				err = fmt.Errorf("%s: %s %s format error: installMode is docker, imageRepo.internal.dataDir can not be empty", errInfo, fieldName, fieldValue)
				return err
			}
			if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
				err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
				return err
			}
		}

		arr := strings.Split(ic.ImageRepo.Internal.Version, ".")
		if len(arr) != 3 {
			fieldName = "imageRepo.internal.version"
			fieldValue = ic.ImageRepo.Internal.Version
			err = fmt.Errorf("%s: %s %s format error: should like v2.4.0", errInfo, fieldName, fieldValue)
			return err
		}
		arr[2] = "0"
		ic.ImageRepo.Internal.VersionBig = strings.Join(arr, ".")

		if ic.ImageRepo.Internal.Password == "" {
			ic.ImageRepo.Internal.Password = RandomString(16, false, "=")
		}
		if ic.ImageRepo.Internal.RegistryPassword == "" {
			ic.ImageRepo.Internal.RegistryPassword = RandomString(16, false, "")
		}
		bs, _ := bcrypt.GenerateFromPassword([]byte(ic.ImageRepo.Internal.RegistryPassword), 10)
		ic.ImageRepo.Internal.RegistryHtpasswd = string(bs)
	}

	fieldName = "hostIP"
	fieldValue = ic.HostIP
	err = ValidateIpAddress(fieldValue)
	if err != nil {
		err = fmt.Errorf("%s: %s %s format error: %s", errInfo, fieldName, fieldValue, err.Error())
		return err
	}
	if fieldValue == "127.0.0.1" || fieldValue == "localhost" {
		err = fmt.Errorf("%s: %s %s format error: can not be 127.0.0.1 or localhost", errInfo, fieldName, fieldValue)
		return err
	}

	var count int
	if ic.Kubernetes.PvConfigLocal.LocalPath != "" {
		count = count + 1
	}
	if len(ic.Kubernetes.PvConfigCephfs.CephMonitors) > 0 {
		count = count + 1
	}
	if ic.Kubernetes.PvConfigNfs.NfsServer != "" {
		count = count + 1
	}
	if count != 1 {
		err = fmt.Errorf("%s: kubernetes.pvConfigLocal/pvConfigNfs/pvConfigCephfs must set one only", errInfo)
		return err
	}

	if ic.Kubernetes.PvConfigLocal.LocalPath != "" {
		if !strings.HasPrefix(ic.Kubernetes.PvConfigLocal.LocalPath, "/") {
			fieldName = "kubernetes.pvConfigLocal.localPath"
			fieldValue = ic.Kubernetes.PvConfigLocal.LocalPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}
	if len(ic.Kubernetes.PvConfigCephfs.CephMonitors) > 0 {
		for _, monitor := range ic.Kubernetes.PvConfigCephfs.CephMonitors {
			fieldName = "kubernetes.pvConfigCephfs.cephMonitors"
			fieldValue = monitor
			arr := strings.Split(monitor, ":")
			if len(arr) != 2 {
				err = fmt.Errorf("%s: %s %s format error: should like 192.168.0.1:6789", errInfo, fieldName, fieldValue)
				return err
			}
			_, err = strconv.Atoi(arr[1])
			if err != nil {
				err = fmt.Errorf("%s: %s %s format error: should like 192.168.0.1:6789", errInfo, fieldName, fieldValue)
				return err
			}
		}
		if ic.Kubernetes.PvConfigCephfs.CephSecret == "" {
			fieldName = "kubernetes.pvConfigCephfs.cephSecret"
			fieldValue = ic.Kubernetes.PvConfigCephfs.CephSecret
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if ic.Kubernetes.PvConfigCephfs.CephUser == "" {
			fieldName = "kubernetes.pvConfigCephfs.cephUser"
			fieldValue = ic.Kubernetes.PvConfigCephfs.CephUser
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if !strings.HasPrefix(ic.Kubernetes.PvConfigCephfs.CephPath, "/") {
			fieldName = "kubernetes.pvConfigCephfs.cephPath"
			fieldValue = ic.Kubernetes.PvConfigCephfs.CephPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}

	if ic.Kubernetes.PvConfigNfs.NfsServer != "" {
		if !strings.HasPrefix(ic.Kubernetes.PvConfigNfs.NfsPath, "/") {
			fieldName = "kubernetes.pvConfigNfs.nfsPath"
			fieldValue = ic.Kubernetes.PvConfigNfs.NfsPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}

	if ic.Dory.Openldap.Password == "" {
		ic.Dory.Openldap.Password = RandomString(16, false, "=")
	}
	if ic.Dory.Redis.Password == "" {
		ic.Dory.Redis.Password = RandomString(16, false, "=")
	}
	if ic.Dory.Mongo.Password == "" {
		ic.Dory.Mongo.Password = RandomString(16, false, "=")
	}

	return err
}

func (ic *InstallConfig) UnmarshalMapValues() (map[string]interface{}, error) {
	var err error
	errInfo := fmt.Sprintf("unmarshal install config to map error")

	bs, _ := yaml.Marshal(ic)
	vals := map[string]interface{}{}
	err = yaml.Unmarshal(bs, &vals)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return vals, err
	}
	imageRepoInternal := true
	imageRepoDomainName := ic.ImageRepo.Internal.DomainName
	imageRepoUsername := "admin"
	imageRepoPassword := ic.ImageRepo.Internal.Password
	imageRepoEmail := "admin@example.com"
	imageRepoIp := ic.HostIP
	if ic.ImageRepo.Internal.DomainName == "" {
		imageRepoInternal = false
		imageRepoDomainName = ic.ImageRepo.External.Url
		imageRepoUsername = ic.ImageRepo.External.Username
		imageRepoPassword = ic.ImageRepo.External.Password
		imageRepoEmail = ic.ImageRepo.External.Email
		imageRepoIp = ic.ImageRepo.External.Ip
	}
	vals["imageRepoInternal"] = imageRepoInternal
	vals["imageRepoDomainName"] = imageRepoDomainName
	vals["imageRepoUsername"] = imageRepoUsername
	vals["imageRepoPassword"] = imageRepoPassword
	vals["imageRepoEmail"] = imageRepoEmail
	vals["imageRepoIp"] = imageRepoIp

	artifactRepoInternal := true
	artifactRepoPortHub := ic.Dory.ArtifactRepo.Internal.PortHub
	artifactRepoPortGcr := ic.Dory.ArtifactRepo.Internal.PortGcr
	artifactRepoPortQuay := ic.Dory.ArtifactRepo.Internal.PortQuay
	artifactRepoUsername := "admin"
	artifactRepoPassword := "Nexus_Pwd_321"
	artifactRepoPublicUser := "public-user"
	artifactRepoPublicPassword := "public-user"
	artifactRepoPublicEmail := "public-user@139.com"
	artifactRepoIp := ic.HostIP
	if ic.Dory.ArtifactRepo.Internal.Image == "" {
		artifactRepoInternal = false
		artifactRepoPortHub = ic.Dory.ArtifactRepo.External.PortHub
		artifactRepoPortGcr = ic.Dory.ArtifactRepo.External.PortGcr
		artifactRepoPortQuay = ic.Dory.ArtifactRepo.External.PortQuay
		artifactRepoUsername = ic.Dory.ArtifactRepo.External.Username
		artifactRepoPassword = ic.Dory.ArtifactRepo.External.Password
		artifactRepoPublicUser = ic.Dory.ArtifactRepo.External.PublicUser
		artifactRepoPublicPassword = ic.Dory.ArtifactRepo.External.PublicPassword
		artifactRepoPublicEmail = ic.Dory.ArtifactRepo.External.PublicEmail
		artifactRepoIp = ic.Dory.ArtifactRepo.External.Host
	}
	vals["artifactRepoInternal"] = artifactRepoInternal
	vals["artifactRepoPortHub"] = artifactRepoPortHub
	vals["artifactRepoPortGcr"] = artifactRepoPortGcr
	vals["artifactRepoPortQuay"] = artifactRepoPortQuay
	vals["artifactRepoUsername"] = artifactRepoUsername
	vals["artifactRepoPassword"] = artifactRepoPassword
	vals["artifactRepoPublicUser"] = artifactRepoPublicUser
	vals["artifactRepoPublicPassword"] = artifactRepoPublicPassword
	vals["artifactRepoPublicEmail"] = artifactRepoPublicEmail
	vals["artifactRepoIp"] = artifactRepoIp

	return vals, err
}

func (ic *InstallConfig) HarborQuery(url, method string, param map[string]interface{}) (string, int, error) {
	var err error
	var strJson string
	var statusCode int
	var req *http.Request
	var resp *http.Response
	var bs []byte
	client := &http.Client{
		Timeout: time.Second * time.Duration(time.Second*5),
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	domainName := ic.ImageRepo.Internal.DomainName
	username := "admin"
	password := ic.ImageRepo.Internal.Password
	if ic.ImageRepo.Internal.DomainName == "" {
		domainName = ic.ImageRepo.External.Url
		username = ic.ImageRepo.External.Username
		password = ic.ImageRepo.External.Password
	}

	url = fmt.Sprintf("https://%s%s", domainName, url)

	if len(param) > 0 {
		bs, err = json.Marshal(param)
		if err != nil {
			return strJson, statusCode, err
		}
		req, err = http.NewRequest(method, url, bytes.NewReader(bs))
		if err != nil {
			return strJson, statusCode, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return strJson, statusCode, err
		}
	}

	req.SetBasicAuth(username, password)
	resp, err = client.Do(req)
	if err != nil {
		return strJson, statusCode, err
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return strJson, statusCode, err
	}
	strJson = string(bs)
	return strJson, statusCode, err
}

func (ic *InstallConfig) HarborProjectAdd(projectName string) error {
	var err error
	var statusCode int
	var strJson string

	url := fmt.Sprintf("/api/v2.0/projects")
	param := map[string]interface{}{
		"project_name": projectName,
		"public":       true,
	}
	strJson, statusCode, err = ic.HarborQuery(url, http.MethodPost, param)
	if err != nil {
		return err
	}

	errmsg := fmt.Sprintf("%s %s", gjson.Get(strJson, "errors.0.code").String(), gjson.Get(strJson, "errors.0.message").String())
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		err = fmt.Errorf(errmsg)
		return err
	}

	return err
}

func (ic *InstallConfig) KubernetesQuery(url, method string, param map[string]interface{}) (string, int, error) {
	var err error
	var strJson string
	var statusCode int
	var req *http.Request
	var resp *http.Response
	var bs []byte
	client := &http.Client{
		Timeout: time.Second * time.Duration(time.Second*5),
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url = fmt.Sprintf("https://%s:%d%s", ic.Kubernetes.Host, ic.Kubernetes.Port, url)

	if len(param) > 0 {
		bs, err = json.Marshal(param)
		if err != nil {
			return strJson, statusCode, err
		}
		req, err = http.NewRequest(method, url, bytes.NewReader(bs))
		if err != nil {
			return strJson, statusCode, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return strJson, statusCode, err
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ic.Kubernetes.Token))
	resp, err = client.Do(req)
	if err != nil {
		return strJson, statusCode, err
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return strJson, statusCode, err
	}
	strJson = string(bs)
	return strJson, statusCode, err
}

func (ic *InstallConfig) KubernetesPodsGet(namespace string) ([]KubePod, error) {
	var err error
	var statusCode int
	var strJson string

	pods := []KubePod{}

	url := fmt.Sprintf("/api/v1/namespaces/%s/pods", namespace)
	param := map[string]interface{}{}
	strJson, statusCode, err = ic.KubernetesQuery(url, http.MethodGet, param)
	if err != nil {
		return pods, err
	}

	errmsg := gjson.Get(strJson, "message").String()
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		err = fmt.Errorf(errmsg)
		return pods, err
	}

	var podList KubePodList
	err = json.Unmarshal([]byte(strJson), &podList)
	if err != nil {
		return pods, err
	}

	pods = podList.Items
	return pods, err
}
