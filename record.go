package iolog

import (
	"fmt"
	"time"
)

const (
	Read  = "read"
	Write = "write"
	Close = "close"
)

type Record struct {
	Tag       string
	Start     time.Time
	Stop      time.Time
	Data      []byte
	Interface interface{}
	Error     error
}

func (r *Record) String() string {
	stop := r.Stop
	if stop.IsZero() {
		stop = r.Start
	}
	var iface string
	if r.Interface != nil { // if not any typed value (but typed <nil> allowed)
		switch i := r.Interface.(type) {
		case fmt.Stringer:
			iface = i.String()
		default:
			iface = fmt.Sprintf(" %v", i)
		}
	}
	const tf = "2006-01-02T15:04:05.000-0700"
	return fmt.Sprintf("%s [% X]%s (%s) %s / %s error: %v", r.Tag, r.Data, iface, stop.Sub(r.Start), r.Start.Format(tf), r.Start.Format(tf), r.Error)
}
