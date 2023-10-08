package test

import (
	"aliyun-oss/conf"
	"aliyun-oss/utils"
	"path/filepath"
	"runtime"
	"testing"
)

func TestConf(t *testing.T) {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("error acquiring current file path")
	}

	confDir := filepath.Dir(fileName)
	conf.InitConf(confDir)
	zero := conf.ConfEntry[string]("")

	if len(zero) != 0 {
		utils.Log(utils.ERRO, "should be blank item")
		t.Error("should be blank item")
	}
}
