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
	// H4sIAAAAAAAA/8pIzcnJ11Eozy/KSQEEAAD//zpyq/8MAAAA
	// [104 101 108 108 111 44 32 119 111 114 108 100]
}
