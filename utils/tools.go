package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func ErrorRecover() {
	if err := recover(); err != nil {
		var buf = make([]byte, 1<<10)
		runtime.Stack(buf, false)
		Log(ERRO, fmt.Sprintf("%v\n%s", err, string(buf)))
	}
}

// generated by GPT4
func ProjectRoot() string {
	defer ErrorRecover()

	// 获取当前执行文件的路径
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("Cannot find executable path: %v", err)
	}

	// 遍历从当前目录到根目录的所有目录
	for dir := filepath.Dir(exe); dir != ""; dir = filepath.Dir(dir) {
		// 检查这个目录是否包含 go.mod 文件
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
	}

	// 如果找不到 go.mod 文件，返回空字符串
	return ""
}

func GenNoteHeader(filePath string) (string, error) {
	fileInfo, statErr := os.Stat(filePath)
	if statErr != nil {
		Log(ERRO, statErr.Error())
		return "", statErr
	}
	pathParts := strings.Split(filepath.Dir(filePath), "/")[1:]
	var categories string = ""
	var tags string = ""
	if len(pathParts) > 0 {
		categories = pathParts[0]
		tags = strings.Join(pathParts, ", ")
	}

	header := ""
	header += "---\n"
	header += fmt.Sprintf("title: %s\n", filepath.Base(filePath)[0:len(filepath.Base(filePath))-len(filepath.Ext(filePath))]) // title
	header += fmt.Sprintf("date: %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))                                     // date
	header += "description:  \" \"\n"                                                                                         // description
	header += fmt.Sprintf("categories: %s\n", "["+categories+"]")                                                             // categories
	header += fmt.Sprintf("tags: %s\n", "["+tags+"]")                                                                         // tags
	header += "---\n"

	return header, nil
}

func ReplaceAllLFInFile(filePath string, lfPattern, newStr string) error {
	defer ErrorRecover()

	header, genHeadErr := GenNoteHeader(filePath)
	if genHeadErr != nil {
		Log(ERRO, fmt.Sprintf("failed to generate header for file: %s", filePath))
	}

	fileContent, readFileErr := os.ReadFile(filePath)
	if readFileErr != nil {
		Log(ERRO, readFileErr.Error())
		return readFileErr
	}

	// hide all codeblocks
	codeBlockRe := regexp.MustCompile("(?s)```.*?```\n?")
	codeBlocks := codeBlockRe.FindAllString(string(fileContent), -1)
	modifiedContent := codeBlockRe.ReplaceAllString(string(fileContent), "CODE_BLOCK_PLACE_HOLDER")

	// define regular expression `(?m)([^\s])\n`
	reLineFeed := regexp.MustCompile(lfPattern)
	modifiedContent = reLineFeed.ReplaceAllString(modifiedContent, newStr)

	// recover codeblocks
	for _, codeBlock := range codeBlocks {
		modifiedContent = strings.Replace(modifiedContent, "CODE_BLOCK_PLACE_HOLDER", codeBlock, 1)
	}

	writeFileErr := os.WriteFile(filePath, []byte(header+modifiedContent), 0611)
	if writeFileErr != nil {
		Log(ERRO, writeFileErr)
		return writeFileErr
	}

	return nil
}
