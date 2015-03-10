// +build debug,!testdeb

package logdeb

import "fmt"

func prDeb(fnc string, par ...interface{}) {
	fmt.Println("*D*", "[["+fnc+"]]", par)
}

func prTest(title string, par ...interface{}) {}
