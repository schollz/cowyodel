package main

import "fmt"

func ExampleEncodeDecode() {
	b1 := []byte("hello, world")
	fmt.Println(b1)
	s1, _ := BytesToString(b1)
	fmt.Println(s1)
	s2, _ := StringToByte(s1)
	fmt.Println(s2)
	// Output: [104 101 108 108 111 44 32 119 111 114 108 100]
	// 68656c6c6f2c20776f726c64
	// [104 101 108 108 111 44 32 119 111 114 108 100]
}
