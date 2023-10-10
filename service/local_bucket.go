package service

import (
	"aliyun-oss/utils"
	"os"
	"path/filepath"
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
		if node.IsDir() {
			if _, exist := remoteObjectMap[localObject]; !exist {
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
			if _, exist := remoteObjectMap[localObject]; !exist {
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
