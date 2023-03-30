package util

import "fmt"

func Println(format string, a ...any) {
	fmt.Printf(fmt.Sprintf("%s\n", format), a...)
}
