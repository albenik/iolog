package iolog

import (
	"fmt"
	"sync"
	"time"
)

type item struct {
	rec  *Record
	next *item
}

var itemPool = sync.Pool{
	New: func() interface{} { return new(item) },
}

func newItem(r *Record) *item {
	i := itemPool.Get().(*item)
	i.rec = r
	i.next = nil
	return i
}

func (i *item) free() {
	i.rec = nil
	i.next = nil
	itemPool.Put(i)
}

type Record struct {
	Tag   string
	Start time.Time
	Stop  time.Time
	Data  interface{}
	Error error
}

var recordPool = sync.Pool{
	New: func() interface{} { return new(Record) },
}

func newRecord(t string, s, f time.Time, e error) *Record {
	r := recordPool.Get().(*Record)
	r.Tag = t
	r.Start = s
	r.Stop = f
	r.Error = e
	r.Data = nil
	return r
}

func (r *Record) free() {
	r.Data = nil
	r.Error = nil
	recordPool.Put(r)
}

func (r *Record) String() string {
	stop := r.Stop
	if stop.IsZero() {
		stop = r.Start
	}
	var datastr string
	switch d := r.Data.(type) {
	case []byte:
		datastr = fmt.Sprintf("[% X]", d)
	case fmt.Stringer:
		datastr = d.String()
	default:
		datastr = fmt.Sprintf("%+v", d)
	}
	const tf = "2006-01-02T15:04:05.000-0700"
	return fmt.Sprintf("%s %s (%s) %s / %s error: %v", r.Tag, datastr, stop.Sub(r.Start), r.Start.Format(tf), r.Start.Format(tf), r.Error)
}
