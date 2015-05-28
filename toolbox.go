package main

import "strings"

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isIPV4(ip string) bool {
	if strings.Contains(ip, ".") {
		return true
	}
	return false
}
