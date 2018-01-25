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

type testConfig struct {
	data  []*iolog.Record
	regex *regexp.Regexp
}

func TestLog(t *testing.T) {
	cfg := []testConfig{
		{
			data: []*iolog.Record{
				{
					Operation: iolog.Read,
					Start:     time.Date(2018, 1, 1, 12, 30, 50, 0, time.Local),
					Stop:      time.Date(2018, 1, 1, 12, 31, 01, 22, time.Local),
					Data:      []byte{1, 2, 3, 4, 5},
					Error:     nil,
				},
			},
			regex: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","iolog":\[{"op":"read","start":1514799050,"stop":1514799061,"dur":11.000000022,"data":"01 02 03 04 05","error":""}]}`),
		}, {
			data: []*iolog.Record{
				{
					Operation: iolog.Write,
					Start:     time.Date(2018, 1, 1, 12, 30, 50, 0, time.Local),
					Stop:      time.Date(2018, 1, 1, 12, 31, 01, 22, time.Local),
					Data:      []byte{1, 2, 3, 4, 5},
					Error:     nil,
				},
			},
			regex: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","iolog":\[{"op":"write","start":1514799050,"stop":1514799061,"dur":11.000000022,"data":"01 02 03 04 05","error":""}]}`),
		}, {
			data: []*iolog.Record{
				{
					Operation: iolog.Close,
					Start:     time.Date(2018, 1, 1, 12, 30, 50, 0, time.Local),
					Stop:      time.Date(2018, 1, 1, 12, 30, 50, 5, time.Local),
					Data:      nil,
					Error:     nil,
				},
			},
			regex: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","iolog":\[{"op":"close","start":1514799050,"stop":1514799050,"dur":0.000000005,"data":"","error":""}]}`),
		},
		{
			data: []*iolog.Record{
				{
					Operation: iolog.Read,
					Start:     time.Date(2018, 1, 1, 12, 30, 50, 0, time.Local),
					Stop:      time.Date(2018, 1, 1, 12, 30, 50, 0, time.Local),
					Data:      []byte{9},
					Error:     io.EOF,
				},
			},
			regex: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","iolog":\[{"op":"read","start":1514799050,"stop":1514799050,"dur":0,"data":"09","error":"EOF"}]}`),
		},
	}
	out := new(bytes.Buffer)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.Lock(zapcore.AddSync(out)), zap.DebugLevel)
	logger := zap.New(core)

	for _, c := range cfg {
		out.Reset()
		logger.Info("Test ok", zap.Array("iolog", iozap.Log(c.data)))
		logger.Sync()
		assert.Regexp(t, c.regex, out.String())
	}
}
