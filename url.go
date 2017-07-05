package main

import (
	"github.com/goware/urlx"
)

func parseURL(urlstring string) (server string, user string, password string, err error) {
	u, err := urlx.Parse(urlstring)
	if err != nil {
		return
	}
	server = u.Scheme + "://" + u.Host
	if u.User != nil {
		user = u.User.Username()
		password, _ = u.User.Password()
	}
	return
}
