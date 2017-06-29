package main

import "testing"

func TestMain(t *testing.T) {
	err := run()
	if err != nil {
		t.Error(err)
	}
}
