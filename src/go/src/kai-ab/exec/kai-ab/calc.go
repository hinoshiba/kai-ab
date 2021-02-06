package main

import (
	"os"
	"fmt"
	"sort"
	"path/filepath"
	"strconv"
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
		r_name := e_name[len(e_name)-4:len(e_name)]

		r, err := do_calc(r_name, srcs)
		if err != nil {
			return err
		}

		if err := do_export(e_path, r); err != nil {
			return err
		}
	}
	return nil
}

func do_calc(name string, srcs []string) (*Report, error) {
	report := NewReport(name)
	for _, src := range srcs {
		r, err := calcMonth(src)
		if err != nil {
			return nil, err
		}

		if err := report.Invite(r); err != nil {
			return nil, err
		}
	}
	return report, nil
}

func do_export(path string, ry *Report) error {
	fpath := path + ".md"

	header := "|項目|01|02|03|04|05|06|07|08|09|10|11|12|\n"
	header += "|:---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n"

	r_in := "|収入|"
	r_dec := "|支出|"
	r_sum := "|黒字値|"

	keys := ry.Keys()
	r_ks := make(map[string]string)


	rms := ry.Childs()

	if len(rms) < 1 {
		return fmt.Errorf("'%s' haven't report", ry.Title())
	}
	h_rm := rms[0]
	m, err := strconv.Atoi(h_rm.Title()[4:])
	if err != nil {
		return err
	}
	h_padval := ""

	for i := m; i > 1; i-- {
		h_padval += "-|"
	}
	r_in  += h_padval
	r_dec += h_padval
	r_sum += h_padval

	f_list := "## 計算元ファイル\n"
	for _, rm := range rms {
		r_in += fmt.Sprintf("%s|", nfmt(rm.IncSum()))
		r_dec += fmt.Sprintf("%s|", nfmt(rm.DecSum()))
		r_sum += fmt.Sprintf("%s|", nfmt(rm.IncSum() + rm.DecSum()))

		for _, k := range keys {
			line, ok := r_ks[k]
			if !ok {
				line = fmt.Sprintf("|%s|%s", k, h_padval)
			}
			val, ok := rm.GetDetail(k)
			if !ok {
				line += "-|"
			}
			line += fmt.Sprintf("%v|", nfmt(val))
			r_ks[k] = line
		}

		f_list += "\n### " + rm.Title() + "\n"
		for _, csv_r := range rm.Childs() {

			rel, err := filepath.Rel(filepath.Dir(fpath), csv_r.Title())
			if err != nil {
				return err
			}
			name := filepath.Base(rel)
			f_list += fmt.Sprintf("* [%s](%s)\n    * 収入(%s) - 支出(%s) = 黒字値(%s)\n",
								name, rel, nfmt(csv_r.IncSum()), nfmt(csv_r.DecSum()),
								nfmt(csv_r.IncSum() + csv_r.DecSum()))
		}
	}

	padval := ""
	pdding_size := 12 - len(rms)
	for i := 0; i < pdding_size; i++ {
		padval += "-|"
	}

	r_str := ry.Title() + "\n===\n\n"
	r_str += header
	r_str += r_in + padval + "\n"
	r_str += r_dec + padval + "\n"
	r_str += r_sum + padval + "\n"
	r_str += "## 詳細\n"
	r_str += header
	for _, line := range r_ks {
		r_str += line + padval + "\n"
	}
	r_str += f_list

	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(r_str); err != nil {
		return err
	}

	return nil
}

func getYears(csv_dir string) (map[string][]string, error) {
	fnames, err := DoLs(csv_dir)
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

func calcMonth(target string) (*Report, error) {
	report_name := target[len(target)-6:len(target)]
	in_dir := filepath.Join(target, PATH_CSV_IN)
	out_dir := filepath.Join(target, PATH_CSV_OUT)

	in_fnames, err := DoLs(in_dir)
	if err != nil {
		return nil, err
	}
	out_fnames, err := DoLs(out_dir)
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

		in_repo := NewReport(in_fpath)
		for _, row := range in_c.Rows() {
			if err := in_repo.Add(row.Category(), row.Size() * SYMBOL_IN); err != nil {
				return nil, err
			}
		}
		if err := report.Invite(in_repo); err != nil {
			return nil, err
		}
	}
	for _, out_fname := range out_fnames {
		out_fpath := filepath.Join(out_dir, out_fname)
		out_c, err := csv.Open(out_fpath, nil)
		if err != nil {
			return nil, err
		}
		defer out_c.Close()

		out_repo := NewReport(out_fpath)
		for _, row := range out_c.Rows() {
			if err := out_repo.Add(row.Category(), row.Size() * SYMBOL_OUT); err != nil {
				return nil, err
			}
		}
		if err := report.Invite(out_repo); err != nil {
			return nil, err
		}
	}

	return report, nil
}

type Report struct {
	title   string
	cl_vals map[string]int64
	inc_sum int64
	dec_sum int64

	childs  []*Report
}

func NewReport(title string) *Report {
	return &Report{
		inc_sum: 0,
		dec_sum: 0,
		title: title,
		cl_vals: make(map[string]int64),
		childs: make([]*Report, 0),
	}
}

func (self *Report) Title() string {
	return self.title
}

func (self *Report) Add(category string, size int64) error {
	return self.add(category, size)
}

func (self *Report) Invite(r *Report) error {
	for r_category, r_val := range r.cl_vals {
		if err := self.add(r_category, r_val); err != nil {
			return err
		}
	}
	self.childs = append(self.childs, r)
	return nil
}

func (self *Report) Childs() []*Report {
	sort.SliceStable(self.childs, func(i, j int) bool { return self.childs[i].Title() < self.childs[j].Title()})
	return self.childs
}

func (self *Report) Keys() []string {
	keys := []string{}
	for key, _ := range self.cl_vals {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool { return keys[i] < keys[j]})
	return keys
}

func (self *Report) GetDetail(key string) (int64, bool) {
	val, ok := self.cl_vals[key]
	if !ok {
		return 0, false
	}
	return val, true
}

func (self *Report) IncSum() int64 {
	return self.inc_sum
}

func (self *Report) DecSum() int64 {
	return self.dec_sum
}

func (self *Report) add(category string, size int64) error {
	if category == "" {
		return fmt.Errorf("cannot append value when empty category")
	}
	sum, ok := self.cl_vals[category]
	if !ok {
		sum = 0
	}
	sum += size
	self.cl_vals[category] = sum

	if 0 <= size {
		self.inc_sum += size
	} else {
		self.dec_sum += size
	}
	return nil
}

func nfmt(num int64) string {
	s := fmt.Sprintf("%d", num)
	fnum := ""
	pos := 0
	for i := len(s) - 1; i >= 0; i-- {
		if pos > 2 && pos % 3 == 0 {
			fnum = fmt.Sprintf(",%s", fnum)
		}
		fnum = fmt.Sprintf("%s%s", string(s[i]), fnum)
		pos++
	}
	return fnum
}
