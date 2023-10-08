package conf

import (
	"aliyun-oss/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var confMap map[string]interface{}

func init() {
	defer utils.ErrorRecover()

	if confMap == nil {
		utils.Log(utils.INFO, "scanning configuration jsons...")
		confMap = make(map[string]interface{}, 0)
	}

	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("error acquiring current conf.go file path")
	}

	confDir := filepath.Dir(fileName)

	dirItems, dirReadErr := os.ReadDir(confDir)

	if dirReadErr != nil {
		panic(dirReadErr.Error())
	}

	type confItem struct {
		key string
		val interface{}
	}

	itemChan := make(chan confItem)

	var wg sync.WaitGroup
	for _, item := range dirItems {
		if item.IsDir() {
			continue
		}

		if !strings.HasSuffix(item.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(confDir, item.Name())

		wg.Add(1)
		go func(filePath string, ch chan confItem) {
			defer utils.ErrorRecover()
			defer wg.Done()

			jsonData, fileReadErr := os.ReadFile(filePath)
			if fileReadErr != nil {
				panic(fileReadErr.Error())
			}

			var tmpConf map[string]interface{}
			unmarshalErr := json.Unmarshal([]byte(jsonData), &tmpConf)
			if unmarshalErr != nil {
				panic(unmarshalErr.Error())
			}

			for key, val := range tmpConf {
				ch <- confItem{key, val}
			}
			utils.Log(utils.INFO, fmt.Sprintf("load config: %s", filePath))
		}(filePath, itemChan)
	}
	// create a routine to accept values  from  channel
	// or else it will wait until the channel is closed
	go func() {
		for item := range itemChan {
			confMap[item.key] = item.val
		}
	}()
	wg.Wait()
	close(itemChan) // don't forget to close the channel
}

func ConfEntryKeys() []string {
	keys := []string{}
	for key := range confMap {
		keys = append(keys, key)
	}

	return keys
}

func ConfEntry[T any](key string) T {
	if _, exist := confMap[key]; !exist {
		var zero T
		return zero
	}
	return confMap[key].(T)
}
