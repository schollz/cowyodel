package main

import (
	"encoding/hex"
)

func BytesToString(b []byte) (converted string, err error) {
	converted = hex.EncodeToString(b)
	return
}

func StringToByte(s string) (converted []byte, err error) {
	converted, err = hex.DecodeString(s)
	return
}
