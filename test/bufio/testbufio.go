package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {

	rd := bufio.NewReader(strings.NewReader("I am a handsome boy, you are a pretty girl,"))
	line, _ := rd.ReadSlice(',')
	fmt.Println(string(line))
	newLine, _ := rd.ReadSlice(',')
	fmt.Println(string(newLine))
	fmt.Println(string(line))
}
