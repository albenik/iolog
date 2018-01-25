package iozap

import (
	"fmt"

	"go.uber.org/zap/zapcore"

	"github.com/albenik/iolog"
)

func record(r *iolog.Record) zapcore.ObjectMarshaler {
	return zapcore.ObjectMarshalerFunc(func(obj zapcore.ObjectEncoder) error {
		obj.AddString("op", string(r.Operation))
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

func Log(log []*iolog.Record) zapcore.ArrayMarshaler {
	return zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
		for _, r := range log {
			arr.AppendObject(record(r))
		}
		return nil
	})
}

// Keep for backward compatibility
func IOLog(log []*iolog.Record) zapcore.ArrayMarshaler {
	return Log(log)
}
