package main

import (
	"fmt"
	"path/filepath"
)

import (
	"kai-ab/csv"
	"kai-ab/conf"
)

type Filter interface {
	Do(string) error
}

func cmd_autofilter(path string) error {
	af, err := NewAutoFilter(CurDir)
	if err != nil {
		return err
	}

	csv_dir := filepath.Join(CurDir, PATH_CSV_BOTH)
	mdirs, err := DoLs(csv_dir)
	if err != nil {
		return err
	}
	for _, mdir := range mdirs {
		path = filepath.Join(csv_dir, mdir)
		if err := do_filter(path, af); err != nil {
			return err
		}
	}
	return nil
}

func cmd_mfilter(path string) error {
	return nil
}

func do_filter(path string, f Filter) error {
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

		if err := f.Do(in_fpath); err != nil {
			return err
		}
	}

	for _, out_fname := range out_fnames {
		out_fpath := filepath.Join(out_dir, out_fname)

		if err := f.Do(out_fpath); err != nil {
			return err
		}
	}

	return nil
}

type AutoFilter struct {
	findex map[string]string
}

func NewAutoFilter(c_path string) (*AutoFilter, error) {
	fs_path := filepath.Join(c_path, PATH_ETC_FILTERS)
	fs, err := conf.LoadFilters(fs_path)
	if err != nil {
		return nil, err
	}

	findex := make(map[string]string)
	for _, f := range fs {
		if _, ok := findex[f.Name]; ok {
			return nil, fmt.Errorf("due name at filters.yaml. '%s'", f.Name)
		}

		findex[f.Name] = f.Category
	}

	return &AutoFilter{findex: findex}, nil
}

func (self *AutoFilter) Do(path string) error {
	c, err := csv.Open(path, nil)
	if err != nil {
		return err
	}

	rows := c.Rows()
	changed := false
	for i, row := range rows {
		if row.Category() != "" {
			continue
		}

		category, ok := self.findex[row.Name()]
		if !ok {
			continue
		}

		row.SetCategory(category)
		rows[i] = row
		changed = true
	}

	if !changed {
		if err := c.Close(); err != nil {
			return err
		}
		return nil
	}
	c.UpdateRows(rows)
	if err := c.CloseWithSave(); err != nil {
		return err
	}
	return nil
}
