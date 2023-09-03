package main

import (
	"os"
	"fmt"
	"path/filepath"
	"bufio"
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
	csv_dir := filepath.Join(CurDir, PATH_CSV_BOTH)
	mdirs, err := DoLs(csv_dir)
	if err != nil {
		return err
	}

	mf := NewManualFilter()
	for _, mdir := range mdirs {
		path = filepath.Join(csv_dir, mdir)
		if err := do_filter(path, mf); err != nil {
			return err
		}
	}
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

type ManualFilter struct {
	findex map[string]string
}

func NewManualFilter() *ManualFilter {
	return &ManualFilter{findex: make(map[string]string)}
}

func (self *ManualFilter) Do(path string) error {
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

		fmt.Printf("[ %s, %s, %v | Category: %s | Memo: %s ] ",
				row.DateString(), row.Name(), row.Size(), row.Category(), row.Memo())
		if row.Category() != "" {
			if str := getstr("Update?[y/N(any key)]"); str != "y" {
				continue
			}
		}

		var category string
		for {
			category = getstr("enter CategoryName > ")
			if category == "" {
				fmt.Println("Error: cannot set empty to CategoryName.")
				continue
			}

			if ans := getstr("Are you sure?[y/N(any)]"); ans != "y" {
				continue
			}

			break
		}
		var memo string
		for {
			memo = getstr("enter Memo > ")
			if ans := getstr("Are you sure?[y/N(any)]"); ans != "y" {
				continue
			}

			break
		}

		row.SetCategory(category)
		row.SetMemo(memo)
		rows[i] = row
		fmt.Printf("Updated !!! [ %s, %s, %v | Category: %s | Memo: %s ]\n",
				row.DateString(), row.Name(), row.Size(), row.Category(), row.Memo())
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

func getstr(msg string) string {
	fmt.Printf(msg)
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	return s.Text()
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
