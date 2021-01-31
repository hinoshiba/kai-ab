package main

import (
	"fmt"
	"path/filepath"
)

import (
	"kai-ab/csv"
)

func cmd_autofilter(path string) error {
	fs_path := filepath.Join(CurDir, PATH_ETC_FILTERS)
	fs, err := LoadFilters(fs_path)
	if err != nil {
		return err
	}
	findex := make(map[string]string)
	for _, f := range fs {
		if _, ok := findex[f.Name]; ok {
			return fmt.Errorf("due name at filters.yaml. '%s'", f.Name)
		}

		findex[f.Name] = f.Category
	}

	csv_dir := filepath.Join(CurDir, PATH_CSV_BOTH)
	mdirs, err := DoLs(csv_dir)
	if err != nil {
		return err
	}
	for _, mdir := range mdirs {
		path = filepath.Join(csv_dir, mdir)
		if err := do_filter(path, findex); err != nil {
			return err
		}
	}
	return nil
}

func do_filter(path string, findex map[string]string) error {
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
		changed := false

		in_fpath := filepath.Join(in_dir, in_fname)
		in_c, err := csv.Open(in_fpath, nil)
		if err != nil {
			return err
		}

		rows := in_c.Rows()
		for i, row := range rows {
			if row.Category() != "" {
				continue
			}

			category, ok := findex[row.Name()]
			if !ok {
				continue
			}

			row.SetCategory(category)
			rows[i] = row
			changed = true
		}

		if !changed {
			if err := in_c.Close(); err != nil {
				return err
			}
			continue
		}
		in_c.UpdateRows(rows)
		if err := in_c.CloseWithSave(); err != nil {
			return err
		}
	}

	for _, out_fname := range out_fnames {
		changed := false

		out_fpath := filepath.Join(out_dir, out_fname)
		out_c, err := csv.Open(out_fpath, nil)
		if err != nil {
			return err
		}

		rows := out_c.Rows()
		for i, row := range rows {
			if row.Category() != "" {
				continue
			}

			category, ok := findex[row.Name()]
			if !ok {
				continue
			}

			row.SetCategory(category)
			rows[i] = row
			changed = true
		}

		if !changed {
			if err := out_c.Close(); err != nil {
				return err
			}
			continue
		}
		out_c.UpdateRows(rows)
		if err := out_c.CloseWithSave(); err != nil {
			return err
		}
	}

	return nil
}
