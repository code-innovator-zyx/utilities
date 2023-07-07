package open

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
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
Label1:
	for _, projects := range core.Editors {
		for project, editor := range projects {
			if path.Base(project) == projectName {
				matches = map[string]string{project: editor}
				break Label1
			}
			if strings.Contains(path.Base(project), projectName) {
				matches[project] = editor
			}
		}
	}
	if len(matches) == 0 {
		fmt.Printf("%s match nothing\n", projectName)
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
	mapping := make(map[string]string, len(projects))
	index := 1
	for p, _ := range projects {
		fmt.Println("              index     	project")
		fmt.Printf("              %d    		%s\n", index, p)
		index += 1
		mapping[strconv.Itoa(index)] = p
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("请选择需要打开的项目索引：")
	indexStr, _ := reader.ReadString('\n')
	indexStr = strings.TrimSpace(indexStr)

	err := exec.Command(projects[mapping[indexStr]], mapping[indexStr]).Run()
	if err != nil {
		panic(err)
	}

	return true
}

func help() {
	fmt.Println("qopen open helps to open the projects ")
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
