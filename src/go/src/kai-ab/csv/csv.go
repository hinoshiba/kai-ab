package csv

import (
	"io"
	"os"
	"fmt"
	"sync"
	"encoding/csv"
)

type Csv struct {
	fpath  string
	fh     *os.File

	header []string
	index  map[string]int

	rows   []*Row

	opt    *Options
	mtx    *sync.Mutex
}

type Options struct {
	Header bool
}

func Open(path string, opt *Options) (*Csv, error) {
	if opt == nil {
		opt = &Options{
			Header: true,
		}
	}

	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(fh)

	var header []string
	var rows   []*Row

	for {
		raw, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if opt.Header {
			if header == nil {
				header = raw
				continue
			}
		}

		row := NewRow(raw)
		rows = append(rows, row)
	}

	index := make(map[string]int)
	if opt.Header {
		if header != nil {
			for i, k := range header {
				index[k] = i
			}
		}
	}

	return &Csv{
		fpath: path,
		fh: fh,
		header: header,
		index: index,

		rows: rows,

		opt: opt,
		mtx: new(sync.Mutex),
	}, nil
}

func (self *Csv) Save() error {
	self.lock()
	defer self.unlock()

	return self.save()
}

func (self *Csv) CloseWithSave() error {
	self.lock()
	defer self.unlock()

	if err := self.save(); err != nil {
		return err
	}
	return self.close()
}

func (self *Csv) Close() error {
	self.lock()
	defer self.unlock()

	return self.close()
}

func (self *Csv) close() error {
	if err := self.fh.Close(); err != nil {
		return err
	}

	self.fh = nil
	self.rows = nil
	self.header = nil
	self.index = nil
	return nil
}

func (self *Csv) save() error {
	if err := self.fh.Truncate(0); err != nil {
		return err
	}

	w := csv.NewWriter(self.fh)
	if self.opt.Header {
		if err := w.Write(self.header); err != nil {
			return err
		}
	}

	for _, row := range self.rows {
		if err := w.Write(row.Raw()); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

func (self *Csv) Header() ([]string, error) {
	self.lock()
	defer self.unlock()

	if !self.opt.Header {
		return nil, fmt.Errorf("does not mode of load the header.")
	}
	if self.header == nil {
		return nil, fmt.Errorf("header is not defined.")
	}
	return self.header, nil
}

func (self *Csv) Path() string {
	self.lock()
	defer self.unlock()

	return self.fpath
}

func (self *Csv) GetIndexId(key string) (int, error) {
	self.lock()
	defer self.unlock()

	if !self.opt.Header {
		return -1, fmt.Errorf("does not mode load of header.")
	}
	if self.index == nil {
		return -1, fmt.Errorf("index is not defined.")
	}

	id, ok := self.index[key]
	if !ok {
		return -1, fmt.Errorf("undefined key, '%s'", key)
	}
	return id, nil
}

func (self *Csv) Rows() ([]*Row, error) {
	self.lock()
	defer self.unlock()

	if self.rows == nil {
		return nil, fmt.Errorf("rows is not defined.")
	}
	return self.rows, nil
}

func (self *Csv) UpdateRows(rows []*Row) error {
	self.lock()
	defer self.unlock()

	if rows != nil {
		return fmt.Errorf("cannot set nil pointer")
	}
	self.rows = rows
	return nil
}

func (self *Csv) UpdateHeader(header []string) error {
	self.lock()
	defer self.unlock()

	if header != nil {
		return fmt.Errorf("cannot set nil pointer")
	}
	if !self.opt.Header {
		return fmt.Errorf("does not mode of load the header.")
	}

	index := make(map[string]int)
	for i, k := range header {
		index[k] = i
	}

	self.index = index
	self.header = header
	return nil
}

func (self *Csv) lock() {
	self.mtx.Lock()
}

func (self *Csv) unlock() {
	self.mtx.Unlock()
}

type Row struct {
	vals []string
}

func NewRow(raw []string) *Row {
	return &Row{
		vals: raw,
	}
}

func (self *Row) Raw() []string {
	return self.vals
}
