package iozap

import (
	"fmt"

	"github.com/albenik/iolog"
	"go.uber.org/zap/zapcore"
)

type objectMarshalerFunc func(zapcore.ObjectEncoder) error

func (fn objectMarshalerFunc) MarshalLogObject(obj zapcore.ObjectEncoder) error {
	return fn(obj)
}

func marshalLogRecord(r *iolog.Record) zapcore.ObjectMarshaler {
	return objectMarshalerFunc(func(obj zapcore.ObjectEncoder) error {
		switch r.Operation {
		case iolog.Read:
			obj.AddString("op", "read")
		case iolog.Write:
			obj.AddString("op", "write")
		case iolog.Close:
			obj.AddString("op", "close")
		default:
			obj.AddString("op", "unknown")
		}

		obj.AddTime("start", r.Start)
		stop := r.Stop
		if stop.IsZero() {
			stop = r.Start
		}
		obj.AddTime("stop", stop)
		obj.AddDuration("dur", stop.Sub(r.Start))
		obj.AddString("data", fmt.Sprintf("% X", r.Data))
		if r.Error == nil {
			obj.AddString("error", "")
		} else {
			obj.AddString("error", r.Error.Error())
		}
		return nil
	})
}

type IOLog []*iolog.Record

func (l IOLog) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for _, r := range l {
		arr.AppendObject(marshalLogRecord(r))
	}
	return nil
}
