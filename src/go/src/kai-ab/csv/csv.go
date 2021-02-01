package csv

import (
	"io"
	"os"
	"fmt"
	"sync"
	"time"
	"strconv"
	"encoding/csv"

	_ "time/tzdata"
)

const (
	FMT_TIME string = "2006/01/02"
)

var Header []string = []string{"date", "name", "size", "category", "memo"}

type Csv struct {
	fpath  string
	fh     *os.File

	header []string

	rows   []*Row

	opt    *Options
	mtx    *sync.Mutex
}

type Options struct {
	Header   bool
	Timezone *time.Location
}

func Open(path string, opt *Options) (*Csv, error) {
	if opt == nil {
		jst, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			return nil, err
		}

		opt = &Options{
			Header: true,
			Timezone: jst,
		}
	}

	fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(fh)
	r.FieldsPerRecord = -1

	var header []string
	var rows   []*Row

	i := 0
	for {
		raw, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		i++
		if opt.Header {
			if header == nil {
				header = Header
				continue
			}
		}

		row, err := NewRow(raw, opt.Timezone)
		if err != nil {
			return nil, fmt.Errorf("cannot load line '%v'. %s", i, err)
		}
		rows = append(rows, row)
	}

	return &Csv{
		fpath: path,
		fh: fh,
		header: header,

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
	return nil
}

func (self *Csv) save() error {
	if err := self.fh.Truncate(0); err != nil {
		return err
	}
	if _, err := self.fh.Seek(0,0); err != nil {
		return err
	}

	if self.header == nil {
		self.header = Header
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

func (self *Csv) Rows() []*Row {
	self.lock()
	defer self.unlock()

	return self.rows
}

func (self *Csv) UpdateRows(rows []*Row) error {
	self.lock()
	defer self.unlock()

	if rows == nil {
		return fmt.Errorf("cannot set nil pointer")
	}
	self.rows = rows
	return nil
}

func (self *Csv) lock() {
	self.mtx.Lock()
}

func (self *Csv) unlock() {
	self.mtx.Unlock()
}

type Row struct {
	date  time.Time
	name  string
	size  int64
	category string
	memo  string
}

func CreateRow(date time.Time, name string, size int64, category string, memo string) *Row {
	return &Row{
		date: date,
		name: name,
		size: size,
		category: category,
		memo: memo,
	}
}

func NewRow(raw []string, tz *time.Location) (*Row, error) {
	if len(raw) < 3 {
		return nil, fmt.Errorf("lack column less than 3.")
	}

	date, err := time.ParseInLocation(FMT_TIME, raw[0], tz)
	if err != nil {
		return nil, err
	}
	size, err := strconv.ParseInt(raw[2], 10, 64)
	if err != nil {
		return nil, err
	}

	var category string
	if len(raw) >= 4 {
		category = raw[3]
	}
	var memo string
	if len(raw) >= 5 {
		memo = raw[4]
	}

	return &Row{
		date: date,
		name: raw[1],
		size: size,
		category: category,
		memo: memo,
	}, nil
}

func (self *Row) Date() time.Time {
	return self.date
}

func (self *Row) DateString() string {
	return self.date.Format(FMT_TIME)
}

func (self *Row) Name() string {
	return self.name
}

func (self *Row) Size() int64 {
	return self.size
}

func (self *Row) Category() string {
	return self.category
}

func (self *Row) Memo() string {
	return self.memo
}

func (self *Row) SetCategory(name string) {
	self.category = name
}

func (self *Row) SetMemo(memo string) {
	self.memo = memo
}

func (self *Row) Raw() []string {
	return self.body2raw()
}

func (self *Row) body2raw() []string {
	raw := make([]string, 5, 5)

	raw[0] = self.date.Format(FMT_TIME)
	raw[1] = self.name
	raw[2] = strconv.FormatInt(self.size, 10)
	raw[3] = self.category
	raw[4] = self.memo

	return raw
}
