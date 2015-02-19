// +build debug

package logdeb

import "fmt"

func prDeb(fnc string, par ...interface{}) {
	fmt.Println("*D*", "[["+fnc+"]]", par)
}
