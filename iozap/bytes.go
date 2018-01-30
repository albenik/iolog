package iozap

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

type StringerFunc func() string

func (fn StringerFunc) String() string {
	return fn()
}

func BytesArray(key string, val []byte) zapcore.Field {
	return zapcore.Field{
		Key:  key,
		Type: zapcore.StringerType,
		Interface: StringerFunc(func() string {
			return fmt.Sprintf("% X", val)
		}),
	}
}
