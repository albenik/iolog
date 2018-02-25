package iolog

import (
	"fmt"
	"time"
)

type Record struct {
	Tag   string
	Start time.Time
	Stop  time.Time
	Data  interface{}
	Error error
}

func (r *Record) String() string {
	stop := r.Stop
	if stop.IsZero() {
		stop = r.Start
	}
	var data string
	if r.Data != nil { // if not any typed value (but typed <nil> allowed)
		switch d := r.Data.(type) {
		case []byte:
			data = fmt.Sprintf("[% X]", d)
		case fmt.Stringer:
			data = d.String()
		default:
			data = fmt.Sprintf("%+v", d)
		}
	}
	const tf = "2006-01-02T15:04:05.000-0700"
	return fmt.Sprintf("%s %s (%s) %s / %s error: %v", r.Tag, data, stop.Sub(r.Start), r.Start.Format(tf), r.Start.Format(tf), r.Error)
}
