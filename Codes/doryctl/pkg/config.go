package pkg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"
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

	fieldName = "dory.gitRepo.imageDB"
	fieldValue = ic.Dory.GitRepo.ImageDB
	if ic.Dory.GitRepo.Type == "gitea" && ic.Dory.GitRepo.ImageDB == "" {
		err = fmt.Errorf("%s: %s %s format error: gitea imageDB can not be empty", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "dory.artifactRepo.type"
	fieldValue = ic.Dory.ArtifactRepo.Type
	if ic.Dory.ArtifactRepo.Type != "nexus" {
		err = fmt.Errorf("%s: %s %s format error: must be nexus", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "imageRepo.type"
	fieldValue = ic.ImageRepo.Type
	if ic.ImageRepo.Type != "harbor" {
		err = fmt.Errorf("%s: %s %s format error: must be harbor", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "imageRepo.namespace"
	fieldValue = ic.ImageRepo.Namespace
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "imageRepo.Namespace"
	fieldValue = ic.ImageRepo.Namespace
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	if ic.InstallMode == "docker" {
		fieldName = "imageRepo.certsDir"
		fieldValue = ic.ImageRepo.CertsDir
		if fieldValue == "" {
			err = fmt.Errorf("%s: %s %s format error: installMode is docker, imageRepo.certsDir can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
			err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
			return err
		}

		fieldName = "imageRepo.dataDir"
		fieldValue = ic.ImageRepo.DataDir
		if fieldValue == "" {
			err = fmt.Errorf("%s: %s %s format error: installMode is docker, imageRepo.dataDir can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
			err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
			return err
		}
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

	arr := strings.Split(ic.ImageRepo.Version, ".")
	if len(arr) != 3 {
		fieldName = "imageRepo.version"
		fieldValue = ic.ImageRepo.Version
		err = fmt.Errorf("%s: %s %s format error: should like v2.4.0", errInfo, fieldName, fieldValue)
		return err
	}
	arr[2] = "0"
	ic.ImageRepo.VersionBig = strings.Join(arr, ".")

	if ic.Dory.Openldap.Password == "" {
		ic.Dory.Openldap.Password = RandomString(16, false, "=")
	}
	if ic.Dory.Redis.Password == "" {
		ic.Dory.Redis.Password = RandomString(16, false, "=")
	}
	if ic.Dory.Mongo.Password == "" {
		ic.Dory.Mongo.Password = RandomString(16, false, "=")
	}
	if ic.ImageRepo.Password == "" {
		ic.ImageRepo.Password = RandomString(16, false, "=")
	}
	if ic.ImageRepo.RegistryPassword == "" {
		ic.ImageRepo.RegistryPassword = RandomString(16, false, "")
	}
	bs, _ := bcrypt.GenerateFromPassword([]byte(ic.ImageRepo.RegistryPassword), 10)
	ic.ImageRepo.RegistryHtpasswd = string(bs)

	return err
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

	url = fmt.Sprintf("https://%s%s", ic.ImageRepo.DomainName, url)

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

	req.SetBasicAuth("admin", ic.ImageRepo.Password)
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
