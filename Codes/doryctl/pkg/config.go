package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

func (idc *InstallDockerConfig) VerifyInstallDockerConfig() error {
	var err error
	errInfo := fmt.Sprintf("verify install docker config error")

	var fieldName, fieldValue string

	fieldName = "rootDir"
	fieldValue = idc.RootDir
	if !strings.HasPrefix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
		return err
	}
	if strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "doryDir"
	fieldValue = idc.DoryDir
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "dory.gitRepo.type"
	fieldValue = idc.Dory.GitRepo.Type
	if idc.Dory.GitRepo.Type != "gitea" && idc.Dory.GitRepo.Type != "gitlab" {
		err = fmt.Errorf("%s: %s %s format error: must be gitea or gitlab", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "dory.gitRepo.imageDB"
	fieldValue = idc.Dory.GitRepo.ImageDB
	if idc.Dory.GitRepo.Type == "gitea" && idc.Dory.GitRepo.ImageDB == "" {
		err = fmt.Errorf("%s: %s %s format error: gitea imageDB can not be empty", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "dory.artifactRepo.type"
	fieldValue = idc.Dory.ArtifactRepo.Type
	if idc.Dory.ArtifactRepo.Type != "nexus" {
		err = fmt.Errorf("%s: %s %s format error: must be nexus", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "imageRepo.type"
	fieldValue = idc.ImageRepo.Type
	if idc.ImageRepo.Type != "harbor" {
		err = fmt.Errorf("%s: %s %s format error: must be harbor", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "imageRepoDir"
	fieldValue = idc.ImageRepoDir
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "imageRepoDir.certsDir"
	fieldValue = idc.ImageRepo.CertsDir
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "imageRepoDir.dataDir"
	fieldValue = idc.ImageRepo.DataDir
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "hostIP"
	fieldValue = idc.HostIP
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
	if idc.Kubernetes.PvConfigLocal.LocalPath != "" {
		count = count + 1
	}
	if len(idc.Kubernetes.PvConfigCephfs.CephMonitors) > 0 {
		count = count + 1
	}
	if idc.Kubernetes.PvConfigNfs.NfsServer != "" {
		count = count + 1
	}
	if count != 1 {
		err = fmt.Errorf("%s: kubernetes.pvConfigLocal/pvConfigNfs/pvConfigCephfs must set one only", errInfo)
		return err
	}

	if idc.Kubernetes.PvConfigLocal.LocalPath != "" {
		if !strings.HasPrefix(idc.Kubernetes.PvConfigLocal.LocalPath, "/") {
			fieldName = "kubernetes.pvConfigLocal.localPath"
			fieldValue = idc.Kubernetes.PvConfigLocal.LocalPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}
	if len(idc.Kubernetes.PvConfigCephfs.CephMonitors) > 0 {
		for _, monitor := range idc.Kubernetes.PvConfigCephfs.CephMonitors {
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
		if idc.Kubernetes.PvConfigCephfs.CephSecret == "" {
			fieldName = "kubernetes.pvConfigCephfs.cephSecret"
			fieldValue = idc.Kubernetes.PvConfigCephfs.CephSecret
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if idc.Kubernetes.PvConfigCephfs.CephUser == "" {
			fieldName = "kubernetes.pvConfigCephfs.cephUser"
			fieldValue = idc.Kubernetes.PvConfigCephfs.CephUser
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if !strings.HasPrefix(idc.Kubernetes.PvConfigCephfs.CephPath, "/") {
			fieldName = "kubernetes.pvConfigCephfs.cephPath"
			fieldValue = idc.Kubernetes.PvConfigCephfs.CephPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}

	if idc.Kubernetes.PvConfigNfs.NfsServer != "" {
		if !strings.HasPrefix(idc.Kubernetes.PvConfigNfs.NfsPath, "/") {
			fieldName = "kubernetes.pvConfigNfs.nfsPath"
			fieldValue = idc.Kubernetes.PvConfigNfs.NfsPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}

	if idc.Dory.Openldap.Password == "" {
		idc.Dory.Openldap.Password = RandomString(16, false, "=")
	}
	if idc.Dory.Redis.Password == "" {
		idc.Dory.Redis.Password = RandomString(16, false, "=")
	}
	if idc.Dory.Mongo.Password == "" {
		idc.Dory.Mongo.Password = RandomString(16, false, "=")
	}
	if idc.ImageRepo.Password == "" {
		idc.ImageRepo.Password = RandomString(16, false, "=")
	}

	return err
}
