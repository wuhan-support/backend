package main

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	Cookie string `yaml:"cookie"`
	Documents struct {
		Accommodations string `yaml:"accommodations"`
		Platforms string `yaml:"platforms"`
	} `yaml:"documents"`
}