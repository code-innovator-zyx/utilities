package main

import (
	"flag"
	"fmt"
	"os"
	"utilities/qopen/file_share"
	"utilities/qopen/list"
	"utilities/qopen/open"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2023/6/29 10:20
* @Package:	快速使用相应编辑器启动电脑上的项目
 */

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		help()
	}
	switch args[0] {
	case "open":
		open.QuickOpen(args[1:]...)
	case "list":
		list.ShowProjects(args[1:]...)
	case "share":
		file_share.Share(args...)
	default:
		help()
	}

}

func help() {
	fmt.Println("qopen is a tool that quickly opens projects using your favorite editor")
	fmt.Println("")
	fmt.Println("Usage")
	fmt.Println("		qopen <command> [arguments]")
	fmt.Println("")
	fmt.Println("The command are:")
	fmt.Println("")
	fmt.Println("          open [project]	open project")
	fmt.Println("          list    list projects")
	os.Exit(0)
}
