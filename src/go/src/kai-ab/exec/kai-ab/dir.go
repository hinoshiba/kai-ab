package main

import (
	"os"
	"fmt"
	"path/filepath"
)

const (
	PATH_ETC_TEMPLATE string = "./etc/template.tml"
	PATH_ETC_ABBASE string = "./etc/accountbook.md.base"
	PATH_ETC_FILTER string = "./etc/filter.tml"
	PATH_ENV string = ".kai-ab.env"
	PATH_REPORT string = "./var/report/"
	PATH_CSV_BOTH string = "./var/csv/"
	PATH_CSV_IN string = "./in/"
	PATH_CSV_OUT string = "./out/"
)

func isKaiabDir(path string) bool {
	p_env := filepath.Join(path, PATH_ENV)
	_, err := os.Stat(p_env)
	return err == nil
}

func cmd_init(path string) error {
	if isKaiabDir(path) {
		return fmt.Errorf("already init")
	}

	p_env := filepath.Join(path, PATH_ENV)
	f_env, err := os.Create(p_env)
	if err != nil {
		return err
	}
	defer f_env.Close()

	_, err = f_env.Write([]byte(VERSION_TEMPLATE))
	if err != nil {
		return err
	}

	p_etc := filepath.Join(path, "etc")
	if err := os.Mkdir(p_etc, 0755); err != nil {
		return err
	}
	p_template := filepath.Join(path, PATH_ETC_TEMPLATE)
	f_t, err := os.Create(p_template)
	if err != nil {
		return err
	}
	f_t.Close()

	p_abbase := filepath.Join(path, PATH_ETC_ABBASE)
	f_ab, err := os.Create(p_abbase)
	if err != nil {
		return err
	}
	f_ab.Close()

	p_filter := filepath.Join(path, PATH_ETC_FILTER)
	f_f, err := os.Create(p_filter)
	if err != nil {
		return err
	}
	f_f.Close()

	p_report := filepath.Join(path, PATH_REPORT)
	if err := os.MkdirAll(p_report, 0755); err != nil {
		return err
	}
	p_csv := filepath.Join(path, PATH_CSV_BOTH)
	if err := os.MkdirAll(p_csv, 0755); err != nil {
		return err
	}


	return nil
}
