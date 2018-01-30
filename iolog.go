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
}

func New(cap int) *IOLog {
	return &IOLog{records: make([]*Record, 0, cap), cap: cap}
}

func (l *IOLog) LogIO(tag string, iofn IOFunc, data []byte) (int, error) {
	rec := &Record{Tag: tag, Start: time.Now()}
	n, err := iofn(data)
	rec.Stop = time.Now()
	rec.Error = err
	if n > 0 {
		rec.Data = make([]byte, n)
		copy(rec.Data, data)
	}
	l.records = append(l.records, rec)
	return n, err
}

func (l *IOLog) LogAny(tag string, fn func(rec *Record) error) error {
	rec := &Record{Tag: tag, Start: time.Now()}
	err := fn(rec)
	l.records = append(l.records, rec)
	return err
}

func (l *IOLog) ClearLog() {
	if cap(l.records) == l.cap {
		l.records = l.records[:0]
	} else {
		l.records = make([]*Record, l.cap)
	}
}

func (l *IOLog) Len() int {
	return len(l.records)
}

func (l *IOLog) Records() []*Record {
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
