package main

import (
	"io/ioutil"
	"path/filepath"
)

import (
	"gopkg.in/yaml.v2"
)

type Template struct {
	In  []*Entry `yaml:"in"`
	Out []*Entry `yaml:"out"`
}

type Entry struct {
	FileName string  `yaml:"fname"`
	Name     string  `yaml:"name"`
	Category string  `yaml:"category"`
	Price    int64   `yaml:"price"`
	Memo     string  `yaml:"memo"`
}

func LoadTemplate(path string) (*Template, error) {
	c_path := filepath.Clean(path)
	yval, err := ioutil.ReadFile(c_path)
	if err != nil {
		return nil, err
	}

	var t Template
	if err := yaml.Unmarshal(yval, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

type Filter struct {
	Name     string `yaml:"name"`
	Category string `yaml:"category"`
}

func LoadFilters(path string) ([]*Filter, error) {
	c_path := filepath.Clean(path)
	yval, err := ioutil.ReadFile(c_path)
	if err != nil {
		return nil, err
	}

	var f []*Filter = []*Filter{}
	if err := yaml.Unmarshal(yval, &f); err != nil {
		return nil, err
	}
	return f, nil
}
