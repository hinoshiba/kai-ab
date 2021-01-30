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

const (
	SYMBOL_IN int64 = 1
	SYMBOL_OUT int64 = -1
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
	log.Println("srcs: %s", srcs)
	report_name := e_path[len(e_path)-4:len(e_path)]

	if err := createDir(e_path); err != nil {
		return err
	}

	report := NewReport(report_name)
	for _, src := range srcs {
		r, err := calcMonth(src)
		if err != nil {
			return err
		}

		report.Merge(r)
	}
	log.Println(report)
	return nil
}

func calcMonth(target string) (*Report, error) {
	report_name := target[len(target)-6:len(target)]
	in_dir := filepath.Join(target, PATH_CSV_IN)
	out_dir := filepath.Join(target, PATH_CSV_OUT)

	in_fnames, err := do_ls(in_dir)
	if err != nil {
		return nil, err
	}
	out_fnames, err := do_ls(out_dir)
	if err != nil {
		return nil, err
	}

	report := NewReport(report_name)
	for _, in_fname := range in_fnames {
		in_fpath := filepath.Join(in_dir, in_fname)
		in_c, err := csv.Open(in_fpath, nil)
		if err != nil {
			return nil, err
		}
		defer in_c.Close()

		for _, row := range in_c.Rows() {
			if err := report.Add(row.Class(), row.Size() * SYMBOL_IN); err != nil {
				return nil, err
			}
		}
	}
	for _, out_fname := range out_fnames {
		out_fpath := filepath.Join(out_dir, out_fname)
		out_c, err := csv.Open(out_fpath, nil)
		if err != nil {
			return nil, err
		}
		defer out_c.Close()

		for _, row := range out_c.Rows() {
			if err := report.Add(row.Class(), row.Size() * SYMBOL_OUT); err != nil {
				return nil, err
			}
		}
	}

	return report, nil
}

func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}
	return nil
}

type Report struct {
	title   string
	cl_vals map[string]int64
	inc_sum int64
	dec_sum int64
}

func NewReport(title string) *Report {
	return &Report{
		inc_sum: 0,
		dec_sum: 0,
		title: title,
		cl_vals: make(map[string]int64),
	}
}

func (self *Report) Add(class string, size int64) error {
	return self.add(class, size)
}

func (self *Report) Merge(r *Report) {
	for r_class, r_val := range r.cl_vals {
		self.add(r_class, r_val)
	}
}

func (self *Report) add(class string, size int64) error {
	if class == "" {
		return fmt.Errorf("cannot append value when empty class")
	}
	sum, ok := self.cl_vals[class]
	if !ok {
		sum = 0
	}
	sum += size
	self.cl_vals[class] = sum

	if 0 <= size {
		self.inc_sum += size
	} else {
		self.dec_sum += size
	}
	return nil
}
