package pkg

import (
	"fmt"
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
