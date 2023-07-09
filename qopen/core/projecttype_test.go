package core

import (
	"fmt"
	"testing"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2023/7/7 15:15
* @Package:
 */

func Test_Get_Editor(t *testing.T) {
	fmt.Println(getEditorAndLanguageFromDir("/Users/zouyuxi/workspace/client"))
}
