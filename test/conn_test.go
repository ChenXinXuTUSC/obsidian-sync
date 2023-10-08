package test

import (
	"aliyun-oss/conf"
	"aliyun-oss/utils"
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func TestOssBucketConnection(t *testing.T) {
	ossEndPoint := conf.ConfEntry[string]("endPoint")
	if len(ossEndPoint) == 0 {
		utils.Log(utils.ERRO, fmt.Sprintf("entry [%s] found", "endPoint"))
		t.Error("entry not found")
	}
	ossAccessKeyId := conf.ConfEntry[string]("accessKeyId")
	if len(ossAccessKeyId) == 0 {
		utils.Log(utils.ERRO, fmt.Sprintf("entry [%s] found", "accessKeyId"))
		t.Error("entry not found")
	}
	ossAccessSecret := conf.ConfEntry[string]("accessKeySecret")
	if len(ossAccessSecret) == 0 {
		utils.Log(utils.ERRO, fmt.Sprintf("entry [%s] found", "accessKeySecret"))
		t.Error("entry not found")
	}

	_, ossClientInitErr := oss.New(
		ossEndPoint,
		ossAccessKeyId,
		ossAccessSecret,
	)

	if ossClientInitErr != nil {
		utils.Log(utils.ERRO, ossClientInitErr.Error())
		t.Error(ossClientInitErr)
	}
}
