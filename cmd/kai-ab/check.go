
package main

import (
	"fmt"
	"path/filepath"
)

import (
	"github.com/hinoshiba/kai-ab/csv"
)

func cmd_check() error {
	csv_dir := filepath.Join(CurDir, PATH_CSV_BOTH)
	mdirs, err := DoLs(csv_dir)
	if err != nil {
		return err
	}
	for _, mdir := range mdirs {
		path := filepath.Join(csv_dir, mdir)
		if err := checkMonthDir(path); err != nil {
			return err
		}
	}
	return nil
}

func checkMonthDir(path string) error {
	in_dir := filepath.Join(path, PATH_CSV_IN)
	out_dir := filepath.Join(path, PATH_CSV_OUT)

	in_fnames, err := DoLs(in_dir)
	if err != nil {
		return err
	}
	out_fnames, err := DoLs(out_dir)
	if err != nil {
		return err
	}

	for _, in_fname := range in_fnames {
		in_fpath := filepath.Join(in_dir, in_fname)

		if err := check(in_fpath); err != nil {
			return err
		}
	}

	for _, out_fname := range out_fnames {
		out_fpath := filepath.Join(out_dir, out_fname)

		if err := check(out_fpath); err != nil {
			return err
		}
	}

	return nil
}

func check(csv_path string) error {
	c, err := csv.Open(csv_path, nil)
	if err != nil {
		return err
	}

	rows := c.Rows()
	for i, row := range rows {
		header_record := 1
		line_num := i + 1 + header_record

		if row.Category() == "" {
			fmt.Printf("%s:%v Not set category.\n", csv_path, line_num)
		}
		if row.Name() == "" {
			fmt.Printf("%s:%v Not set name.\n", csv_path, line_num)
		}
	}
	return nil
}
