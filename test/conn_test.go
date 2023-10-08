package test

import (
	"aliyun-oss/conf"
	"aliyun-oss/utils"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func TestOssBucketConnection(t *testing.T) {
	ossEndPoint := conf.ConfEntry[string]("endPoint")
	ossAccessKeyId := conf.ConfEntry[string]("accessKeyId")
	ossAccessSecret := conf.ConfEntry[string]("accessKeySecret")

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
