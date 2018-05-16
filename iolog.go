package iolog

import (
	"bytes"
	"errors"
	"time"
)

type IOFunc func([]byte) (int, error)

type IOLog struct {
	records []*Record
	active  bool
}

func New() *IOLog {
	return &IOLog{}
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
			p := make([]byte, n)
			copy(p, data)
			rec.Data = p
		}
		l.records = append(l.records, rec)
	}
	return n, err
}

func (l *IOLog) LogAny(tag string, fn func() (interface{}, error)) error {
	start := time.Now()
	data, err := fn()
	if l.active {
		rec := &Record{
			Tag:   tag,
			Start: start,
			Stop:  time.Now(),
			Error: err,
		}
		switch src := data.(type) {
		case []byte:
			if len(src) > 0 {
				p := make([]byte, len(src))
				copy(p, src)
				rec.Data = p
			}
		default:
			rec.Data = data
		}
		l.records = append(l.records, rec)
	}
	return err
}

func (l *IOLog) Start() {
	l.records = make([]*Record, 0, 128)
	l.active = true
}

func (l *IOLog) Stop() []*Record {
	l.active = false
	r := l.records
	l.records = nil
	return r
}

func (l *IOLog) Len() int {
	return len(l.records)
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
