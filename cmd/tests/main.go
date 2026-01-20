package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(len(strings.Split("\r\nabcde\r\nfghij\r\nklmnop\r\n", "\r\n")))
	fmt.Println(len(strings.Split("\r\n", "\r\n")))
	fmt.Println(strings.Index("\r\nabcde", "\r\n"))
	fmt.Println(len("abcde\r\n"))
}
