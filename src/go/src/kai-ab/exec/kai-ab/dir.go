package main

import (
	"os"
	"fmt"
	"path/filepath"
)

func CreateDir(path string) error {
	if !IsExist(path) {
		return os.Mkdir(path, 0755)
	}
	return nil
}

func IsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsKaiabDir(path string) bool {
	p_env := filepath.Join(path, PATH_ENV)
	_, err := os.Stat(p_env)
	return err == nil
}

func cmd_init(path string) error {
	if IsKaiabDir(path) {
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
