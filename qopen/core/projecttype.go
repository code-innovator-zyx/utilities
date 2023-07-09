package core

import (
	"fmt"
	"os"
	"path/filepath"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2023/7/7 14:48
* @Package: 判断根据项目目录获取项目语言类型返回对应打开最合适的编辑器
 */

const (
	GOLAND  = "goland"
	PYCHARM = "charm"
	IDEA    = "idea"
)

const (
	JAVA   = ".java"
	GO     = ".go"
	PYTHON = ".py"
)

// 传入一个文件路径 获取编辑器
func getEditorAndLanguageFromDir(path string) (string, string) {
	var (
		ed       string
		language string
	)
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		switch filepath.Ext(path) {
		case JAVA:
			ed = IDEA
			language = "java"
			return fmt.Errorf("found java file")
		case GO:
			ed = GOLAND
			language = "go"
			return fmt.Errorf("found go file")
		case PYTHON:
			ed = PYCHARM
			language = "py"
			return fmt.Errorf("found python file")
		}
		return nil
	})
	return ed, language
}
