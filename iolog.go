package iolog

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/albenik/goerrors"
)

type Operation string

const (
	Read  = "read"
	Write = "write"
	Close = "close"
)

type Record struct {
	Operation Operation
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
		iface = fmt.Sprintf(" %#v", r.Interface)
	}
	const tf = "2006-01-02T15:04:05.000-0700"
	return fmt.Sprintf("%s [% X]%s (%s) %s / %s error: %v", r.Operation, r.Data, iface, stop.Sub(r.Start), r.Start.Format(tf), r.Start.Format(tf), r.Error)
}

type Wrapper struct {
	reader io.Reader
	writer io.Writer
	closer io.Closer
	log    []*Record
}

func WrapReader(r io.Reader) *Wrapper {
	return &Wrapper{reader: r}
}

func WrapWriter(w io.Writer) *Wrapper {
	return &Wrapper{writer: w}
}

func WrapReadWriter(rw io.ReadWriter) *Wrapper {
	return &Wrapper{reader: rw, writer: rw}
}

func WrapReadCloser(rc io.ReadCloser) *Wrapper {
	return &Wrapper{reader: rc, closer: rc}
}

func WrapWriteCloser(wc io.WriteCloser) *Wrapper {
	return &Wrapper{writer: wc, closer: wc}
}

func WrapReadWriteCloser(rwc io.ReadWriteCloser) *Wrapper {
	return &Wrapper{reader: rwc, writer: rwc, closer: rwc}
}

func (wr *Wrapper) logIO(fn func([]byte) (int, error), op Operation, p []byte) (int, error) {
	rec := &Record{Operation: op, Start: time.Now()}
	n, err := fn(p)
	rec.Stop = time.Now()
	rec.Error = err
	if n > 0 {
		rec.Data = make([]byte, n)
		copy(rec.Data, p)
	}
	wr.log = append(wr.log, rec)
	return n, err
}

func (wr *Wrapper) Read(p []byte) (int, error) {
	if wr.reader == nil {
		err := errors.New("read not possible")
		wr.log = append(wr.log, &Record{Operation: Read, Start: time.Now(), Error: err})
		return 0, err
	}
	return wr.logIO(wr.reader.Read, Read, p)
}

func (wr *Wrapper) Write(p []byte) (int, error) {
	if wr.writer == nil {
		err := errors.New("write not possible")
		wr.log = append(wr.log, &Record{Operation: Write, Start: time.Now(), Error: err})
		return 0, err
	}
	return wr.logIO(wr.writer.Write, Write, p)
}

func (wr *Wrapper) Close() error {
	if wr.closer == nil {
		err := errors.New("close not possible")
		wr.log = append(wr.log, &Record{Operation: Close, Start: time.Now(), Error: err})
		return err
	}
	rec := &Record{Operation: Close, Start: time.Now()}
	err := wr.closer.Close()
	rec.Stop = time.Now()
	rec.Error = err
	wr.log = append(wr.log, rec)
	return err
}

func (wr *Wrapper) AppendLogRecord(r *Record) {
	wr.log = append(wr.log, r)
}

func (wr *Wrapper) LastLogRecord() *Record {
	if len(wr.log) > 0 {
		return wr.log[len(wr.log)-1]
	}
	return nil
}

func (wr *Wrapper) Log() []*Record {
	return wr.log
}

func (wr *Wrapper) ClearLog() {
	wr.log = nil // TODO just reset slice length and keep capacity with configurable maximum allowed value
}

func (wr *Wrapper) String() string {
	var buf bytes.Buffer
	for _, rec := range wr.log {
		buf.WriteString(rec.String())
		buf.WriteRune('\n')
	}
	return buf.String()
}
