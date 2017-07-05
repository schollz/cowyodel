package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/schollz/lumber"
)

var SERVER_STRING string

func init() {
	SERVER_STRING = "localhost:8050"
}

func TestCowyo(t *testing.T) {
	log = lumber.NewConsoleLogger(lumber.TRACE)
	exists, err := pageExists(SERVER_STRING, "alsdkfjalksdfjlaf")
	if err != nil {
		t.Error(err)
	}
	if exists != false {
		t.Error("alsdkfjalksdfjlaf should not exist!")
	}

	err = uploadData(SERVER_STRING, "testpage", "testtext", false, true)
	if err != nil {
		t.Error(err)
	}
	exists, err = pageExists(SERVER_STRING, "testpage")
	if err != nil {
		t.Error(err)
	}
	if exists != true {
		t.Error("testpage should exist!")
	}

	err = downloadData(SERVER_STRING, "testpage", "")
	if err != nil {
		t.Error(err)
	}

	d, err := ioutil.ReadFile("testpage")
	if err != nil {
		t.Error(err)
	}
	if string(d) != "testtext" {
		t.Errorf("Got [%s] instead of testtext", string(d))
	}
	os.Remove("testpage")
}
