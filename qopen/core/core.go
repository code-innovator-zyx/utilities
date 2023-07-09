package core

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2023/6/30 11:03
* @Package:
 */
var projectCharacters = []string{".vscode", ".idea"}

// Dir 自定义目录结构，记录所有子目录
type Dir[T string | struct{}] map[string]T

// Editors 初始化工作目录,获取目录下的项目列表
var (
	Editors map[string]Dir[string]
)

func init() {
	Editors = make(map[string]Dir[string])
	readProjectDirs()
}

type Config struct {
	WorkSpaces    []string `json:"workspaces"`
	DefaultEditor string   `json:"default_editor"`
}

func readProjectDirs() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// 读取文件内容
	content, err := os.ReadFile(path.Join(filepath.Dir(exePath), "/conf/config.yaml"))
	if err != nil {
		panic(err)
	}

	// 将文件内容解析为 map[string]interface{} 类型
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
	for _, workspace := range config.WorkSpaces {
		// 读取目录下最多三级的项目列表
		pickProject(workspace, 1)
	}
}

// 递归获取项目列表
func pickProject(dir string, deep int) {
	if deep >= 3 {
		return
	}
	readDir, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	for _, de := range readDir {
		if !de.IsDir() {
			continue
		}
		for _, character := range projectCharacters {
			targetFile := fmt.Sprintf("%s/%s/%s", dir, de.Name(), character)
			if _, err := os.Stat(targetFile); nil == err {
				projectDir := fmt.Sprintf("%s/%s", dir, de.Name())
				// 获取项目编辑器和对应语言
				editor, language := getEditorAndLanguageFromDir(projectDir)
				if projects, ok := Editors[language]; ok {
					projects[projectDir] = editor
					break
				}
				Editors[language] = Dir[string]{
					projectDir: editor,
				}
				break
			}
		}
		tmp := fmt.Sprintf("%s/%s", dir, de.Name())
		pickProject(tmp, deep+1)
	}
}

// 递归获取项目列表
func findProject(d Dir[string], editor string, dir string, deep int) {
	if deep >= 3 {
		return
	}
	readDir, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	for _, de := range readDir {
		if !de.IsDir() {
			continue
		}
		for _, character := range projectCharacters {
			targetFile := fmt.Sprintf("%s/%s/%s", dir, de.Name(), character)
			if _, err := os.Stat(targetFile); nil == err {
				d[fmt.Sprintf("%s/%s", dir, de.Name())] = editor
				break
			}
		}
		tmp := fmt.Sprintf("%s/%s", dir, de.Name())
		findProject(d, editor, tmp, deep+1)
	}
}
