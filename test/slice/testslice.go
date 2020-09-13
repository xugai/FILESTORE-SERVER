package main

import "fmt"

func main() {

	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s1 := s[0 : 4]
	fmt.Println("s1: ", s1)
	s[0] = 2
	s1 = s[0 : 4]
	fmt.Println("s1: ", s1)
	//s2 := s[4 : 8]
	//fmt.Println("s2: " ,s2)
	//fmt.Println("s1: ", s1)
}
