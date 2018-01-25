package iozap_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/albenik/iolog/iozap"
)

func TestBytesStringer(t *testing.T) {
	testdata := []struct {
		data   []byte
		expect *regexp.Regexp
	}{{
		data:   []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
		expect: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","bytes":"00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F"}`),
	}, {
		data:   []byte{},
		expect: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","bytes":""}`),
	}, {
		data:   nil,
		expect: regexp.MustCompile(`{"level":"info","ts":\d+\.\d+,"msg":"Test ok","bytes":""}`),
	}}
	out := new(bytes.Buffer)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.Lock(zapcore.AddSync(out)), zap.DebugLevel)
	logger := zap.New(core)

	for _, tt := range testdata {
		out.Reset()
		logger.Info("Test ok", zap.Stringer("bytes", iozap.BytearrayStringer(tt.data)))
		logger.Sync()
		assert.Regexp(t, tt.expect, out.String())
	}
}
