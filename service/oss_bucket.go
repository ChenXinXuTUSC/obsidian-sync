package service

import (
	"aliyun-oss/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossClient *oss.Client
var ossBucket *oss.Bucket

func checkOssClientInit() bool {
	return ossClient != nil
}

func checkOssBucketInit() bool {
	return ossBucket != nil
}

func OssClientInit(ossEndPoint, ossAccessKeyId, ossAccessSecret string) error {
	defer utils.ErrorRecover()

	if checkOssClientInit() {
		return nil // already initialized
	}

	client, initErr := oss.New(
		ossEndPoint,
		ossAccessKeyId,
		ossAccessSecret,
	)
	if initErr != nil {
		utils.Log(utils.ERRO, initErr.Error())
		return initErr
	}

	ossClient = client
	return nil
}

func OssBucketInit(bucketName string) error {
	defer utils.ErrorRecover()

	if !checkOssClientInit() {
		utils.Log(utils.WARN, "oss client hasn't been initialized yet")
		return errors.New("oss client hasn't been initialized yet")
	}

	if checkOssBucketInit() {
		return nil // already initialized
	}

	bucket, bucketGetErr := ossClient.Bucket(bucketName)
	if bucketGetErr != nil {
		utils.Log(utils.ERRO, bucketGetErr.Error())
		return bucketGetErr
	}
	ossBucket = bucket

	return nil
}

func OssBucketObjectlist(bucketName string) ([]string, error) {
	defer utils.ErrorRecover()

	if !checkOssClientInit() {
		utils.Log(utils.WARN, "oss client hasn't been initialized yet")
		return nil, errors.New("oss client hasn't been initialized yet")
	}

	if !checkOssBucketInit() {
		ossBucketInitErr := OssBucketInit(bucketName)
		if ossBucketInitErr != nil {
			utils.Log(utils.ERRO, ossBucketInitErr.Error())
			return nil, ossBucketInitErr
		}
	}

	objectList := []string{}
	ossListMarker := ""
	for {
		listRes, listErr := ossBucket.ListObjects(oss.Marker(ossListMarker))
		if listErr != nil {
			utils.Log(utils.ERRO, listErr.Error())
			return nil, listErr
		}

		// contains 100 records by default
		for _, object := range listRes.Objects {
			objectList = append(objectList, object.Key)
		}

		// Flag indicates if all results are returned
		// true if there are still more records
		// false if current list is not truncated which means the end part
		if listRes.IsTruncated {
			ossListMarker = listRes.NextMarker
		} else {
			break
		}
	}

	return objectList, nil
}

func OssBucketObjectListToMap(objectList []string) map[string]struct{} {
	objectMap := make(map[string]struct{})
	for _, objectStr := range objectList {
		objectMap[objectStr] = struct{}{}
	}

	return objectMap
}

func OssBucketObjectDownload(bucketName, object, outDir string) error {
	defer utils.ErrorRecover()

	if strings.HasSuffix(object, "/") {
		// utils.Log(utils.WARN, fmt.Sprintf("%s is a dir not a file", object))
		return fmt.Errorf("%s is a dir not a file", object)
	}

	if !checkOssClientInit() {
		utils.Log(utils.WARN, "oss client hasn't been initialized yet")
		return errors.New("oss client hasn't been initialized yet")
	}

	if !checkOssBucketInit() {
		utils.Log(utils.WARN, "oss bucket hasn't been initialized yet")
		return errors.New("oss bucket hasn't been initialized yet")
	}

	outDir = filepath.Join(outDir, filepath.Dir(object))

	filePath := filepath.Join(outDir, filepath.Base(object))

	if !needToUpdate(object, filePath) {
		return nil
	}

	if _, statErr := os.Stat(outDir); statErr != nil {
		if os.IsNotExist(statErr) {
			errMkdirAll := os.MkdirAll(outDir, 0755)
			if errMkdirAll != nil {
				utils.Log(utils.ERRO, errMkdirAll.Error())
				return errMkdirAll
			}
		}
	}

	// make sure the output directory is exist before
	// file is downloaded to there
	downloadErr := ossBucket.GetObjectToFile(object, filePath)
	if downloadErr != nil {
		utils.Log(utils.ERRO, downloadErr.Error())
		return downloadErr
	}

	utils.Log(utils.INFO, fmt.Sprintf("update object: %s", object))

	return nil
}

func obtainGmtModTimeStr(object string) (string, error) {
	defer utils.ErrorRecover()

	if !checkOssBucketInit() {
		utils.Log(utils.WARN, "bucket has not been initialized yet")
		return "", errors.New("bucket has not been initialized yet")
	}

	objectAttrs, getObjectAttrErr := ossBucket.GetObjectDetailedMeta(object)
	if getObjectAttrErr != nil {
		utils.Log(utils.ERRO, getObjectAttrErr.Error())
		return "", getObjectAttrErr
	}
	// RFC1123 time format, is predefined in golang, can directly
	// be parsed by golang time package
	objectLastModifiedStr := objectAttrs.Get("Last-Modified")

	return objectLastModifiedStr, nil
}

func needToUpdate(object, filePath string) bool {
	defer utils.ErrorRecover()

	fileInfo, statErr := os.Stat(filePath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return true
		}
		utils.Log(utils.WARN, "failed to obtain node status")
		return false
	}

	objectModifiedTime, objModTmeErr := obtainGmtModTimeStr(object)
	if objModTmeErr != nil {
		utils.Log(utils.ERRO, fmt.Sprintf("%s: failed to convert time string", object))
		return false
	}

	// obtain the last modified time
	localLastModified := fileInfo.ModTime()
	remotLastModified, _ := time.Parse(time.RFC1123, objectModifiedTime)

	return remotLastModified.After(localLastModified)
}
