// +build !debug,testdeb

package logdeb

import (
	"fmt"
	"strings"
)

func prDeb(fnc string, par ...interface{}) {}

func prTest(title string, par ...interface{}) {
	var out string
	for i := 0; i < len(par); i++ {
		switch t := par[i].(type) {
		case []byte:
			fmt.Println("QUI")
			out = out + fmt.Sprintf("%v", strings.Trim(strings.Trim(string(par[i].(byte)), "\r"), "\n"))
		case string:
			out = out + fmt.Sprintf("%v", strings.Trim(strings.Trim(par[i].(string), "\r"), "\n"))
		default:
			fmt.Printf("Type: %T", t)
			out = out + fmt.Sprintf("%v", par[i])
		}
	}
	fmt.Println(title, out)
}
