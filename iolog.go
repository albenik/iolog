package iolog

import (
	"bytes"
	"time"

	"github.com/albenik/goerrors"
)

type IOFunc func([]byte) (int, error)

type IOLog struct {
	records []*Record
	cap     int // Initial capacity
	active  bool
}

func New(cap int) *IOLog {
	return &IOLog{records: make([]*Record, 0, cap), cap: cap}
}

func (l *IOLog) LogIO(tag string, iofn IOFunc, data []byte) (int, error) {
	start := time.Now()
	n, err := iofn(data)
	if l.active {
		rec := &Record{
			Tag:   tag,
			Start: start,
			Stop:  time.Now(),
			Error: err,
		}
		if n > 0 {
			rec.Data = make([]byte, n)
			copy(rec.Data, data)
		}
		l.records = append(l.records, rec)
	}
	return n, err
}

func (l *IOLog) LogAny(tag string, fn func() ([]byte, interface{}, error)) error {
	start := time.Now()
	data, iface, err := fn()
	if l.active {
		rec := &Record{
			Tag:       tag,
			Start:     start,
			Stop:      time.Now(),
			Interface: iface,
			Error:     err,
		}
		if len(data) > 0 {
			rec.Data = make([]byte, len(data))
			copy(rec.Data, data)
		}
		l.records = append(l.records, rec)
	}
	return err
}

func (l *IOLog) Len() int {
	return len(l.records)
}

func (l *IOLog) ClearLog() {
	l.active = true
	if cap(l.records) == l.cap {
		l.records = l.records[:0]
	} else {
		l.records = make([]*Record, l.cap)
	}
}

func (l *IOLog) Records() []*Record {
	l.active = false
	return l.records
}

func (l *IOLog) LastRecord() *Record {
	if len(l.records) == 0 {
		now := time.Now()
		l.records = append(l.records, &Record{Tag: "error", Start: now, Error: errors.New("iolog.LastRecord() called for empty log")})
	}
	return l.records[len(l.records)-1]
}

func (l *IOLog) String() string {
	var buf bytes.Buffer
	for _, rec := range l.records {
		buf.WriteString(rec.String())
		buf.WriteRune('\n')
	}
	return buf.String()
}
