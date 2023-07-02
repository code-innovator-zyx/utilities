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
	// 指定编辑器
	for editor, projects := range core.Editors {
		if editor != args[0] {
			continue
		}
		fmt.Printf("Editor %s:\n", editor)
		for name, _ := range projects {
			fmt.Println("              " + name)
		}
	}
}
func help() {
	fmt.Println("qopen l helps to list all the projects ")
	fmt.Println("")
	fmt.Println("Select the project witch you want to open")
	for editor, projects := range core.Editors {
		fmt.Printf("Editor %s:\n", editor)
		for name, _ := range projects {
			fmt.Println("              " + name)
		}
	}
	os.Exit(0)
}
