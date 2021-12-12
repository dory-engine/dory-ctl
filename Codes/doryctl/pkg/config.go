package pkg

import (
	"fmt"
	"strconv"
	"strings"
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

	fieldName = "doryDir"
	fieldValue = ic.DoryDir
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

	fieldName = "imageRepoDir"
	fieldValue = ic.ImageRepoDir
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "imageRepoDir.certsDir"
	fieldValue = ic.ImageRepo.CertsDir
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "imageRepoDir.dataDir"
	fieldValue = ic.ImageRepo.DataDir
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
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

	return err
}
