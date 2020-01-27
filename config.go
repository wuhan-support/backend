package main

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	Cookie string `yaml:"cookie"`
	Documents struct {
		Hotel string `yaml:"hotel"`
	} `yaml:"documents"`
}