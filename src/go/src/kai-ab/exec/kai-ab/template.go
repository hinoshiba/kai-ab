package main

import (
	"fmt"
	"time"
	"path/filepath"
)

import (
	"kai-ab/csv"
	"kai-ab/conf"
)

func cmd_template(date string) error {
	if len(date) < 6 {
		return fmt.Errorf("The number of characters is short less than 6. '%s'", date)
	}
	write_dir := filepath.Join(CurDir, PATH_CSV_BOTH, date)
	if IsExist(write_dir) {
		return fmt.Errorf("target dir is already exist. '%s'", write_dir)
	}

	template_path := filepath.Join(CurDir, PATH_ETC_TEMPLATE)
	template, err := conf.LoadTemplate(template_path)
	if err != nil {
		return err
	}

	date_t, err := time.Parse(PATH_FMT_DATE, date)
	if err != nil {
		return err
	}
	if err := CreateDir(write_dir); err != nil {
		return err
	}

	write_dir_in := filepath.Join(write_dir, PATH_CSV_IN)
	if err := CreateDir(write_dir_in); err != nil {
		return err
	}
	in_bodys := make(map[string][]*csv.Row)
	for _, r := range template.In {
		row := csv.CreateRow(date_t, r.Name, r.Price, r.Category, r.Memo)

		rows, ok := in_bodys[r.FileName]
		if !ok {
			rows = []*csv.Row{}
		}
		rows = append(rows, row)

		in_bodys[r.FileName] = rows
	}
	for fname, rows := range in_bodys {
		path := filepath.Join(write_dir_in, fname)
		c, err := csv.Open(path, nil)
		if err != nil {
			return err
		}

		c.UpdateRows(rows)
		if err := c.CloseWithSave(); err != nil {
			return err
		}
	}

	write_dir_out := filepath.Join(write_dir, PATH_CSV_OUT)
	if err := CreateDir(write_dir_out); err != nil {
		return err
	}
	out_bodys := make(map[string][]*csv.Row)
	for _, r := range template.Out {
		row := csv.CreateRow(date_t, r.Name, r.Price, r.Category, r.Memo)

		rows, ok := out_bodys[r.FileName]
		if !ok {
			rows = []*csv.Row{}
		}
		rows = append(rows, row)

		out_bodys[r.FileName] = rows
	}
	for fname, rows := range out_bodys {
		path := filepath.Join(write_dir_out, fname)
		c, err := csv.Open(path, nil)
		if err != nil {
			return err
		}

		c.UpdateRows(rows)
		if err := c.CloseWithSave(); err != nil {
			return err
		}
	}

	return nil
}
