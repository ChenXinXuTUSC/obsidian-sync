package main

import (
	"aliyun-oss/conf"
	"aliyun-oss/service"
	"aliyun-oss/utils"
	"flag"
	"fmt"
	"path/filepath"
)

var WorkingDirectory string

func init() {
	var confPath string
	flag.StringVar(&confPath, "cf", "./conf", "path to the configuration directory")
	flag.StringVar(&WorkingDirectory, "wd", ".", "path to the working directory")
	flag.Parse()


	conf.InitConf(confPath)
	utils.InitLogFile(conf.ConfEntry[string]("logDir"))
	utils.DumpToFile = true
}

func main() {
	defer utils.ErrorRecover()

	utils.Log(utils.INFO, fmt.Sprintf("working directory has been set to: %s", WorkingDirectory))

	ossClientInitErr := service.OssClientInit(
		conf.ConfEntry[string]("endPoint"),
		conf.ConfEntry[string]("accessKeyId"),
		conf.ConfEntry[string]("accessKeySecret"),
	)
	if ossClientInitErr != nil {
		panic(ossClientInitErr)
	}

	objectList, getObjectListErr := service.OssBucketObjectlist(
		conf.ConfEntry[string]("bucketName"),
	)
	if getObjectListErr != nil {
		panic(getObjectListErr)
	}

	downloadDir := filepath.Join(WorkingDirectory, "data")
	for _, object := range objectList {
		service.OssBucketObjectDownload(
			conf.ConfEntry[string]("bucketName"),
			object,
			downloadDir,
		)
	}
}
