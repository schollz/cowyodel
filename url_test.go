package main

import (
	"testing"
)

func TestURLParse(t *testing.T) {
	s, u, p, err := parseURL("cowyo.com")
	if err != nil {
		t.Error(err)
	}
	if s != "http://cowyo.com" || u != "" || p != "" {
		t.Errorf("Incorrect s,u,p")
		t.Errorf(s)
		t.Errorf(u)
		t.Errorf(p)
	}
	s, u, p, err = parseURL("user:password@cowyo.com")
	if err != nil {
		t.Error(err)
	}
	if s != "http://cowyo.com" || u != "user" || p != "password" {
		t.Errorf("Incorrect s,u,p")
		t.Errorf(s)
		t.Errorf(u)
		t.Errorf(p)
	}
	s, u, p, err = parseURL("https://user:password@cowyo.com")
	if err != nil {
		t.Error(err)
	}
	if s != "https://cowyo.com" || u != "user" || p != "password" {
		t.Errorf("Incorrect s,u,p")
		t.Errorf(s)
		t.Errorf(u)
		t.Errorf(p)
	}
}
