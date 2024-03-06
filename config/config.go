package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)


var Cfg Config


type endpoint struct {
	Path string `json:"path"`
	Permissions map[string]string `json:"permissions"`
}

type service struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Health string `json:"health"`
	Endpoints map[string]endpoint `json:"endpoints"`
	Permissions []string `json:"permissions"`
}

type Config struct {
	Pattern string `json:"pattern"`
	Prefix string `json:"prefix"`
	Services map[string]service `json:"services"`
	Permissions []string `json:"permissions"`
	SecretKey string `json:"secret-key"`
}

func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("File config not loaded")
	}
	json.Unmarshal(data, &Cfg)
	Cfg.Pattern = strings.ReplaceAll(Cfg.Pattern, "{service}", "(?P<service>[a-z-_]+)")
	Cfg.Pattern = strings.ReplaceAll(Cfg.Pattern, "{version}", "(?P<version>[a-z0-9-_.]+)")
	Cfg.Pattern = strings.ReplaceAll(Cfg.Pattern, "{endpoint}", "(?P<endpoint>[a-z-_]+)")
	Cfg.Pattern = strings.ReplaceAll(Cfg.Pattern, "{*}", "(?P<remaining>.*)")
	for sk, sv := range Cfg.Services {
		for ek, ev := range sv.Endpoints {
			if ev.Path == "" {
				ev.Path = fmt.Sprintf("/%s", ek)
			}
			Cfg.Services[sk].Endpoints[ek] = ev
		}
	}
}