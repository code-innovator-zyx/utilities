package image

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2022/11/22 10:05
* @Package:
 */

func Test_Image(t *testing.T) {
	cli := http.DefaultClient
	t.Run("test png", func(t *testing.T) {
		res, err := cli.Get("https://upfile2.asqql.com/upfile/2009pasdfasdfic2009s305985-ts/2014-10/2014101319585958225.gif")
		//res, err := cli.Get("https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fpreview.qiantucdn.com%2F58pic%2F35%2F45%2F89%2F52i58PICkj40bV0wtUWw4_PIC2018.jpg%21w1024_new_small&refer=http%3A%2F%2Fpreview.qiantucdn.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1656151488&t=ee017d7059197f387c029518c133c0f1")
		if nil != err {
			fmt.Println(err)
			return
		}
		b, _ := ioutil.ReadAll(res.Body)
		im, e := DecodeImage(b)
		fmt.Println(e)
		fmt.Printf("%+v", im)
	})
}
