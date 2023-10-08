package main

import (
	"aliyun-oss/conf"
	"aliyun-oss/service"
	"aliyun-oss/utils"
	"flag"
	"fmt"
)

func init() {
	var confPath string
	flag.StringVar(&confPath, "confpath", "./conf", "path to the configuration directory")
	flag.Parse()


	conf.InitConf(confPath)
	utils.InitLogFile(conf.ConfEntry[string]("logDirName"))
	utils.DumpToFile = true
}

func main() {
	defer utils.ErrorRecover()

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

	for _, obj := range objectList {
		fmt.Println(obj)
	}

}
