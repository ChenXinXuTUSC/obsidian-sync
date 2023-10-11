package service

import (
	"aliyun-oss/utils"
	"os"
	"path/filepath"
	"strings"
)

func SyncLocalToRemote(remoteObjectMap map[string]struct{}, targetDir string) error {
	defer utils.ErrorRecover()

	nodes, readDirErr := os.ReadDir(targetDir)
	if readDirErr != nil {
		utils.Log(utils.ERRO, readDirErr.Error())
		return readDirErr
	}

	for _, node := range nodes {
		localObject := filepath.Join(targetDir, node.Name())
		// remove the first part
		checkKey := strings.Join(strings.Split(localObject, "/")[1:], "/")
		if node.IsDir() {
			if _, exist := remoteObjectMap[checkKey + "/"]; !exist {
				utils.Log(utils.DBUG, localObject, "failed test")
				// delete this directory and all files and directories under it
				removeAllErr := os.RemoveAll(localObject)
				if removeAllErr != nil {
					utils.Log(utils.ERRO, removeAllErr.Error())
					return removeAllErr
				}
			} else {
				go SyncLocalToRemote(remoteObjectMap, localObject)
			}
		} else {
			if _, exist := remoteObjectMap[checkKey]; !exist {
				utils.Log(utils.DBUG, localObject, "failed test")
				// delete this file
				removeErr := os.Remove(localObject)
				if removeErr != nil {
					utils.Log(utils.ERRO, removeErr.Error())
					return removeErr
				}
			}
		}
	}

	return nil
}

func ReplaceLineFeed(targetDir string) error {
	defer utils.ErrorRecover()

	nodes, readDirErr := os.ReadDir(targetDir)
	if readDirErr != nil {
		utils.Log(utils.ERRO, readDirErr.Error())
		return readDirErr
	}

	for _, node := range nodes {
		localObject := filepath.Join(targetDir, node.Name())
		if node.IsDir() {
			ReplaceLineFeed(localObject)
		} else {
			if !strings.HasSuffix(node.Name(), ".md") {
				continue
			}
			reErr := utils.ReplaceAllLFInFile(localObject, `(?m)([^\s])\n`, "$1  \n")
			if reErr != nil {
				utils.Log(utils.ERRO, reErr.Error())
				return reErr
			}
		}
	}

	return nil
}
