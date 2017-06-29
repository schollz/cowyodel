package main

import "testing"

func TestNames(t *testing.T) {
	if GetRandomName(1) == GetRandomName(0) {
		t.Errorf("RandomName is not Random!")
	}
}
