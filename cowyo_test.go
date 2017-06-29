package main

import (
	"io/ioutil"
	"os"
	"testing"
)

var SERVER string

func init() {
	SERVER = "http://localhost:8050"
}

func TestCowyo(t *testing.T) {
	exists, err := pageExists(SERVER, "alsdkfjalksdfjlaf")
	if err != nil {
		t.Error(err)
	}
	if exists != false {
		t.Error("alsdkfjalksdfjlaf should not exist!")
	}

	err = uploadData(SERVER, "testpage", "testtext", false, true)
	if err != nil {
		t.Error(err)
	}
	exists, err = pageExists(SERVER, "testpage")
	if err != nil {
		t.Error(err)
	}
	if exists != true {
		t.Error("testpage should exist!")
	}

	err = downloadData(SERVER, "testpage", "")
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
