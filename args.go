package main

import (
	"os"
	"strings"
)


func LoadArgs() map[string]any {
	args := map[string]any{}
	for idx, arg := range os.Args {
		if idx == 0 {
			continue
		}
		splt := strings.Split(arg, "=")
		print(splt)
	}
	return args
}