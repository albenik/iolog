package iolog

import (
	"errors"
	"fmt"
	"time"
)

type IOFunc func([]byte) (int, error)

type IOLog struct {
	active  bool
	first   *item
	current *item
	len     int
	maxlen  int
}

func New(l int) *IOLog {
	if l < 3 {
		panic(errors.New("iolog: maxlen too small"))
	}
	return &IOLog{maxlen: l}
}

func (l *IOLog) append(r *Record) {
	i := newItem(r)

	if l.current == nil {
		l.first = i
	} else {
		l.current.next = i
	}
	l.current = i

	if l.len < l.maxlen {
		l.len++
		return
	}

	if l.len >= l.maxlen {
		first := l.first
		l.first = first.next
		first.free()
	}
}

func (l *IOLog) LogIO(t string, fn IOFunc, p []byte) (int, error) {
	start := time.Now()
	n, err := fn(p)
	if l.active {
		r := newRecord(t, start, time.Now(), err)
		if n > 0 {
			data := make([]byte, n)
			copy(data, p)
			r.Data = data
		}
		l.append(r)
	}
	return n, err
}

func (l *IOLog) LogAny(t string, fn func() (interface{}, error)) error {
	start := time.Now()
	data, err := fn()
	if l.active {
		r := newRecord(t, start, time.Now(), err)
		switch src := data.(type) {
		case []byte:
			if len(src) > 0 {
				p := make([]byte, len(src))
				copy(p, src)
				r.Data = p
			}
		default:
			r.Data = data
		}
		l.append(r)
	}
	return err
}

func (l *IOLog) Start() {
	l.first = nil
	l.current = nil
	l.len = 0
	l.active = true
}

func (l *IOLog) Stop() []*Record {
	l.active = false
	list := make([]*Record, 0, l.len)

	item := l.first
	for item != nil {
		list = append(list, item.rec)
		item = item.next
	}
	return list
}

func (l *IOLog) Len() int {
	return l.len
}

func (l *IOLog) LastRecord() *Record {
	return l.current.rec
}

func (l *IOLog) String() string {
	return fmt.Sprintf("iolog{active:%v, len:%d, maxlen:%d}", l.active, l.len, l.maxlen)
}
