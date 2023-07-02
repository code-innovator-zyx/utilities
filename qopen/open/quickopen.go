package open

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"utilities/qopen/core"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2023/6/30 10:56
* @Package:	快速打开电脑中任意指定项目，未指定编辑器使用默认编辑器，指定编辑器使用指定的编辑器打开
 */

// QuickOpen 打开相应项目
func QuickOpen(args ...string) {
	if len(args) == 0 {
		help()
	}
	var (
		projectName = args[0]
	)
	matches := make(map[string]string)
	// 可能有多个匹配的项目，如果有多个，需要返回让用户重新选择
	for _, projects := range core.Editors {
		for project, editor := range projects {
			if strings.Contains(path.Base(project), projectName) {
				matches[project] = editor
			}
		}
	}
	if len(matches) == 0 {
		help()
	}
	if !matchMany(matches) {
		for p, editor := range matches {
			fmt.Printf("quickly open %s by %s\n", p, editor)
			err := exec.Command(editor, p).Run()
			if err != nil {
				panic(err)
			}
		}
	}

}
func matchMany(projects map[string]string) bool {
	if len(projects) <= 1 {
		return false
	}
	fmt.Println("to many matches:")
	for p, _ := range projects {
		fmt.Println("              " + p)
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("请输入完整项目地址：")
	addr, _ := reader.ReadString('\n')
	addr = strings.TrimSpace(addr)

	err := exec.Command(projects[addr], addr).Run()
	if err != nil {
		panic(err)
	}

	return true
}

func help() {
	fmt.Println("qopen o helps to open the projects ")
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
