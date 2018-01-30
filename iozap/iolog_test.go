package iozap_test

import (
	"bytes"
	"io"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/albenik/iolog"
	"github.com/albenik/iolog/iozap"
)

func TestLog(t *testing.T) {
	testdata := []struct {
		data  []*iolog.Record
		regex *regexp.Regexp
		msg   string
	}{
		{
			data: []*iolog.Record{
				{
					Tag:   iolog.Read,
					Start: time.Date(2018, 1, 1, 12, 30, 01, 11111, time.UTC),
					Stop:  time.Date(2018, 1, 1, 12, 31, 02, 22222, time.UTC),
					Data:  []byte{1, 2, 3, 4, 5},
					Error: nil,
				},
			},
			regex: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","iolog":\[{"op":"read","start":1514809801.000011,"stop":1514809862.0000222,"dur":61.000011111,"data":"01 02 03 04 05","error":""}]}`),
			msg:   "log 1",
		}, {
			data: []*iolog.Record{
				{
					Tag:   iolog.Write,
					Start: time.Date(2018, 1, 1, 12, 30, 01, 11111, time.UTC),
					Stop:  time.Date(2018, 1, 1, 12, 31, 02, 22222, time.UTC),
					Data:  []byte{1, 2, 3, 4, 5},
					Error: nil,
				},
			},
			regex: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","iolog":\[{"op":"write","start":1514809801.000011,"stop":1514809862.0000222,"dur":61.000011111,"data":"01 02 03 04 05","error":""}]}`),
			msg:   "log 2",
		}, {
			data: []*iolog.Record{
				{
					Tag:   iolog.Close,
					Start: time.Date(2018, 1, 1, 12, 30, 50, 11111, time.UTC),
					Stop:  time.Date(2018, 1, 1, 12, 30, 50, 22222, time.UTC),
					Data:  nil,
					Error: nil,
				},
			},
			regex: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","iolog":\[{"op":"close","start":1514809850.000011,"stop":1514809850.0000222,"dur":0.000011111,"data":"","error":""}]}`),
			msg:   "log 3",
		},
		{
			data: []*iolog.Record{
				{
					Tag:   iolog.Read,
					Start: time.Date(2018, 1, 1, 12, 30, 50, 11111, time.UTC),
					Stop:  time.Date(2018, 1, 1, 12, 30, 50, 11111, time.UTC),
					Data:  []byte{9},
					Error: io.EOF,
				},
			},
			regex: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","iolog":\[{"op":"read","start":1514809850.000011,"stop":1514809850.000011,"dur":0,"data":"09","error":"EOF"}]}`),
			msg:   "log 4",
		},
	}
	out := new(bytes.Buffer)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(zapcore.AddSync(out)),
		zap.DebugLevel,
	)
	logger := zap.New(core)

	for _, tt := range testdata {
		out.Reset()
		logger.Info("Test ok", zap.Array("iolog", iozap.Log(tt.data)))
		logger.Sync()
		assert.Regexp(t, tt.regex, out.String(), tt.msg)
	}
}
