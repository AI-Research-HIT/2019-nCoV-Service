package util

import "github.com/ender-wan/ewlog"

func Recover() {
	if r := recover(); r != nil {
		ewlog.Error("Recovered in f ", r)
	}
}
