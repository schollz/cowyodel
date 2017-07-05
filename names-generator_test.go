package main

import (
	"fmt"
	"testing"
)

func TestNames(t *testing.T) {
	fmt.Println(GetRandomName())
	if GetRandomName() == GetRandomName() {
		t.Errorf("RandomName is not Random!")
	}
}
