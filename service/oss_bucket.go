package service

import (
	"aliyun-oss/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossClient *oss.Client

func checkOssClientInit() bool {
	return ossClient != nil
}

func OssClientInit(ossEndPoint, ossAccessKeyId, ossAccessSecret string) error {
	defer utils.ErrorRecover()

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

func OssBucketObjectlist(bucketName string) ([]string, error) {
	defer utils.ErrorRecover()

	if !checkOssClientInit() {
		utils.Log(utils.WARN, "oss client hasn't been initialized yet")
		return nil, errors.New("oss client hasn't been initialized yet")
	}

	bucket, bucketGetErr := ossClient.Bucket(bucketName)
	if bucketGetErr != nil {
		utils.Log(utils.ERRO, bucketGetErr.Error())
		return nil, bucketGetErr
	}

	objectList := []string{}
	ossListMarker := ""
	for {
		listRes, listErr := bucket.ListObjects(oss.Marker(ossListMarker))
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

func OssBucketObjectDownload(bucketName, object, outDir string) error {
	defer utils.ErrorRecover()

	if strings.HasSuffix(object, "/") {
		utils.Log(utils.WARN, fmt.Sprintf("%s is a dir not a file", object))
		return fmt.Errorf("%s is a dir not a file", object)
	}

	if !checkOssClientInit() {
		utils.Log(utils.WARN, "oss client hasn't been initialized yet")
		return errors.New("oss client hasn't been initialized yet")
	}

	bucket, bucketGetErr := ossClient.Bucket(bucketName)
	if bucketGetErr != nil {
		utils.Log(utils.ERRO, bucketGetErr.Error())
		return bucketGetErr
	}


	outDir = filepath.Join(outDir, filepath.Dir(object))
	if _, statErr := os.Stat(outDir); statErr != nil {
		if os.IsNotExist(statErr) {
			errMkdirAll := os.MkdirAll(outDir, 0755)
			if errMkdirAll != nil {
				utils.Log(utils.ERRO, errMkdirAll.Error())
				return errMkdirAll
			}
		}
	}

	filePath := filepath.Join(outDir, filepath.Base(object))

	downloadErr := bucket.GetObjectToFile(object, filePath)
	if downloadErr != nil {
		utils.Log(utils.ERRO, downloadErr.Error())
		return downloadErr
	}

	return nil
}
