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



	targetDir := filepath.Join(WorkingDirectory, conf.ConfEntry[string]("dumpDir"))
	// delete all old directories and files first
	objectRecords := service.OssBucketObjectListToMap(objectList)
	syncLocalErr := service.SyncLocalToRemote(objectRecords, targetDir)
	if syncLocalErr != nil {
		utils.Log(utils.ERRO, syncLocalErr.Error())
		panic(syncLocalErr)
	}

	// update all new directories and files later
	for _, object := range objectList {
		service.OssBucketObjectDownload(
			conf.ConfEntry[string]("bucketName"),
			object,
			targetDir,
		)
	}

	// reformat the line feed as I didn't write
	// use standard line feed for markdown.
	// reErr := service.ReplaceLineFeed(targetDir)
	// if reErr != nil {
	// 	panic(reErr)
	// }
}
