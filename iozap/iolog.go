package iozap

import (
	"fmt"

	"go.uber.org/zap/zapcore"

	"github.com/albenik/iolog"
)

func record(r *iolog.Record) zapcore.ObjectMarshaler {
	return zapcore.ObjectMarshalerFunc(func(obj zapcore.ObjectEncoder) error {
		obj.AddString("op", string(r.Tag))
		obj.AddTime("start", r.Start)
		stop := r.Stop
		if stop.IsZero() {
			stop = r.Start
		}
		obj.AddTime("stop", stop)
		obj.AddDuration("dur", stop.Sub(r.Start))
		if r.Data != nil {
			switch d := r.Data.(type) {
			case []byte:
				obj.AddString("data", fmt.Sprintf("% X", d))
			case fmt.Stringer:
				obj.AddString("data", d.String())
			default:
				obj.AddString("data", fmt.Sprintf("%+v", d))
			}
		}
		if r.Error != nil {
			obj.AddString("error", r.Error.Error())
		}
		return nil
	})
}

func IOLog(key string, log []*iolog.Record) zapcore.Field {
	return zapcore.Field{
		Key:  key,
		Type: zapcore.ArrayMarshalerType,
		Interface: zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
			for _, r := range log {
				arr.AppendObject(record(r))
			}
			return nil
		}),
	}
}
