package test

import (
	"aliyun-oss/conf"
	"aliyun-oss/utils"
	"testing"
)

func TestConf(t *testing.T) {
	zero := conf.ConfEntry[string]("")

	if len(zero) != 0 {
		utils.Log(utils.ERRO, "should be blank item")
		t.Error("should be blank item")
	}
}
