package file_share

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2024/3/4 09:58
* @Package:
 */
var (
	port = flag.String("p", "8080", "Share 文件共享监听端口号")
)

func Share(args ...string) {
	os.Args = args
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir("/Users/zouyuxi/Desktop")))
	fmt.Printf("share file with port [%s]\n", *port)
	e := http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
	fmt.Println(e)
}
