package main

import (
	"aliyun-oss/conf"
	"aliyun-oss/utils"
)

func init() {
	utils.InitLogFile(conf.ConfEntry("logDirName").(string))
	utils.DumpToFile = true
}

func main() {
	utils.Log(utils.INFO, conf.ConfEntryKeys())
}
