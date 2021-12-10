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

	fieldName = "harborDir"
	fieldValue = idc.HarborDir
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "harbor.certsDir"
	fieldValue = idc.Harbor.CertsDir
	if strings.HasPrefix(fieldValue, "/") || strings.HasSuffix(fieldValue, "/") {
		err = fmt.Errorf("%s: %s %s format error: can not start or end with /", errInfo, fieldName, fieldValue)
		return err
	}

	fieldName = "harbor.dataDir"
	fieldValue = idc.Harbor.DataDir
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
	if idc.Dorycore.Kubernetes.PvConfigLocal.LocalPath != "" {
		count = count + 1
	}
	if len(idc.Dorycore.Kubernetes.PvConfigCephfs.CephMonitors) > 0 {
		count = count + 1
	}
	if len(idc.Dorycore.Kubernetes.PvConfigGlusterfs.EndpointIPs) > 0 {
		count = count + 1
	}
	if idc.Dorycore.Kubernetes.PvConfigNfs.NfsServer != "" {
		count = count + 1
	}
	if count != 1 {
		err = fmt.Errorf("%s: dorycore.kubernetes.pvConfigLocal/pvConfigNfs/pvConfigCephfs/pvConfigGlusterfs must set one only", errInfo)
		return err
	}

	if idc.Dorycore.Kubernetes.PvConfigLocal.LocalPath != "" {
		if !strings.HasPrefix(idc.Dorycore.Kubernetes.PvConfigLocal.LocalPath, "/") {
			fieldName = "dorycore.kubernetes.pvConfigLocal.localPath"
			fieldValue = idc.Dorycore.Kubernetes.PvConfigLocal.LocalPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}
	if len(idc.Dorycore.Kubernetes.PvConfigCephfs.CephMonitors) > 0 {
		for _, monitor := range idc.Dorycore.Kubernetes.PvConfigCephfs.CephMonitors {
			fieldName = "dorycore.kubernetes.pvConfigCephfs.cephMonitors"
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
		if idc.Dorycore.Kubernetes.PvConfigCephfs.CephSecret == "" {
			fieldName = "dorycore.kubernetes.pvConfigCephfs.cephSecret"
			fieldValue = idc.Dorycore.Kubernetes.PvConfigCephfs.CephSecret
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if idc.Dorycore.Kubernetes.PvConfigCephfs.CephUser == "" {
			fieldName = "dorycore.kubernetes.pvConfigCephfs.cephUser"
			fieldValue = idc.Dorycore.Kubernetes.PvConfigCephfs.CephUser
			err = fmt.Errorf("%s: %s %s format error: can not be empty", errInfo, fieldName, fieldValue)
			return err
		}
		if !strings.HasPrefix(idc.Dorycore.Kubernetes.PvConfigCephfs.CephPath, "/") {
			fieldName = "dorycore.kubernetes.pvConfigCephfs.cephPath"
			fieldValue = idc.Dorycore.Kubernetes.PvConfigCephfs.CephPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}

	if len(idc.Dorycore.Kubernetes.PvConfigGlusterfs.EndpointIPs) > 0 {
		for _, epi := range idc.Dorycore.Kubernetes.PvConfigGlusterfs.EndpointIPs {
			fieldName = "dorycore.kubernetes.pvConfigGlusterfs.endpointIPs"
			fieldValue = epi
			err = ValidateIpAddress(epi)
			if err != nil {
				err = fmt.Errorf("%s: %s %s format error: %s", errInfo, fieldName, fieldValue, err.Error())
				return err
			}
		}
		if idc.Dorycore.Kubernetes.PvConfigGlusterfs.EndpointPort < 1 {
			fieldName = "dorycore.kubernetes.pvConfigGlusterfs.endpointPort"
			fieldValue = fmt.Sprintf("%d", idc.Dorycore.Kubernetes.PvConfigGlusterfs.EndpointPort)
			err = fmt.Errorf("%s: %s %s format error: must set", errInfo, fieldName, fieldValue)
			return err
		}
		if strings.HasPrefix(idc.Dorycore.Kubernetes.PvConfigGlusterfs.Path, "/") {
			fieldName = "dorycore.kubernetes.pvConfigGlusterfs.path"
			fieldValue = idc.Dorycore.Kubernetes.PvConfigGlusterfs.Path
			err = fmt.Errorf("%s: %s %s format error: can not start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}

	if idc.Dorycore.Kubernetes.PvConfigNfs.NfsServer != "" {
		if !strings.HasPrefix(idc.Dorycore.Kubernetes.PvConfigNfs.NfsPath, "/") {
			fieldName = "dorycore.kubernetes.pvConfigNfs.nfsPath"
			fieldValue = idc.Dorycore.Kubernetes.PvConfigNfs.NfsPath
			err = fmt.Errorf("%s: %s %s format error: must start with /", errInfo, fieldName, fieldValue)
			return err
		}
	}

	if idc.Dory.Gitea.DbPassword == "" {
		idc.Dory.Gitea.DbPassword = RandomString(16, false, "=")
	}
	if idc.Dory.Openldap.AdminPassword == "" {
		idc.Dory.Openldap.AdminPassword = RandomString(16, false, "=")
	}
	if idc.Dory.Redis.Password == "" {
		idc.Dory.Redis.Password = RandomString(16, false, "=")
	}
	if idc.Dory.Mongo.Password == "" {
		idc.Dory.Mongo.Password = RandomString(16, false, "=")
	}
	if idc.Harbor.AdminPassword == "" {
		idc.Harbor.AdminPassword = RandomString(16, false, "=")
	}
	if idc.Harbor.DbPassword == "" {
		idc.Harbor.DbPassword = RandomString(16, false, "=")
	}

	return err
}
