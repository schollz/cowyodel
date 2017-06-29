package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
)

func BytesToString(b []byte) (converted string, err error) {
	var buf bytes.Buffer
	wb64 := base64.NewEncoder(base64.StdEncoding, &buf)
	wgz := gzip.NewWriter(wb64)
	_, err = wgz.Write(b)
	if err != nil {
		return
	}
	wgz.Close()
	wb64.Close()
	converted = string(buf.Bytes())

	return
}

func StringToByte(s string) (converted []byte, err error) {
	var buf bytes.Buffer
	buf.WriteString(s)
	r, err := gzip.NewReader(base64.NewDecoder(base64.StdEncoding, &buf))
	if err != nil {
		return
	}
	converted, err = ioutil.ReadAll(r)
	return
}
