package list

import (
	"fmt"
	"os"
	"utilities/qopen/core"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2023/6/30 11:07
* @Package:	获取本地项目列表 可使用编辑器进行选择 也可自定义编辑器
 */

// ShowProjects 获取项目列表
func ShowProjects(args ...string) {
	if len(args) == 0 {
		help()
	}
	// 指定语言
	for language, projects := range core.Editors {
		if language != args[0] {
			continue
		}
		fmt.Printf("Language %s:\n", language)
		for name, _ := range projects {
			fmt.Println("              " + name)
		}
	}
}
func help() {
	fmt.Println("qopen list helps to list all the projects ")
	fmt.Println("")
	fmt.Println("Current projects list")
	fmt.Println(core.Editors)
	for language, projects := range core.Editors {
		fmt.Printf("Language %s:\n", language)
		for name, _ := range projects {
			fmt.Println("              " + name)
		}
	}
	os.Exit(0)
}
