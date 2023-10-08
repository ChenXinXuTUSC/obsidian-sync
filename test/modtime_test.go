package test

import (
	"aliyun-oss/utils"
	"fmt"
	"os"
	"testing"
)

func TestReadModTime(t *testing.T) {
	fileInfo, err := os.Stat("/home/fredom/workspace/goland/aliyun-bucket/service/oss_bucket.go")
	if err != nil {
		utils.Log(utils.ERRO, err)
		t.Error(err.Error())
	}

	// obtain the last modified time
	lastModified := fileInfo.ModTime()

	utils.Log(utils.INFO, fmt.Sprintf("Last Modified Time of the File: %v\n", lastModified))
}
