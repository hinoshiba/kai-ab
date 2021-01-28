package main

import "log"

import (
	"os"
	"fmt"
	"path/filepath"
)

import (
	"kai-ab/csv"
)

func cmd_calc() error {
	csv_dir := filepath.Join(CurDir, PATH_CSV_BOTH)

	targets, err := getYears(csv_dir)
	if err != nil {
		return err
	}

	for e_name, srcs := range targets {
		e_path := filepath.Join(CurDir, PATH_REPORT, e_name)

		if err := do_calc(srcs, e_path); err != nil {
			return err
		}
	}
	return nil
}

func getYears(csv_dir string) (map[string][]string, error) {
	fnames, err := do_ls(csv_dir)
	if err != nil {
		return nil, err
	}

	years := make(map[string][]string)
	for _, fname := range fnames {
		if len(fname) != 6 {
			err_path := filepath.Join(csv_dir, fname)
			return nil, fmt.Errorf("dirname size error. Need set 6 chars. : '%s'", err_path)
		}

		year := fname[:4]
		paths, ok := years[year]
		if !ok {
			paths = []string{}
		}

		paths = append(paths, filepath.Join(csv_dir, fname))
		years[year] = paths
	}
	return years, nil
}

func do_ls(path string) ([]string, error) {
	d, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	fs, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, f := range fs {
		names = append(names, f.Name())
	}
	return names, nil
}

func do_calc(srcs []string, e_path string) error {
	log.Println("called do_calc, export path is %s", e_path)
	log.Println("srcs: %s", srcs)

	if err := createDir(e_path); err != nil {
		return err
	}

	for _, src := range srcs {
		in_dir := filepath.Join(src, PATH_CSV_IN)
//		out_dir := filepath.Join(src, PATH_CSV_OUT)

		in_fnames, err := do_ls(in_dir)
		if err != nil {
			return err
		}

		for _, in_fname := range in_fnames {
			in_fpath := filepath.Join(in_dir, in_fname)
			in_c, err := csv.Open(in_fpath, nil)
			if err != nil {
				return err
			}
			defer in_c.Close()

			log.Println(in_c.Header())
			log.Println(in_c.Rows())
		}
	}

	return nil
}

func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}
	return nil
}
