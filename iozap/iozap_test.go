package iozap_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/albenik/iolog"
	"github.com/albenik/iolog/iozap"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func Test(t *testing.T) {
	src := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	wr := iolog.WrapReader(bytes.NewReader(src))
	buf := make([]byte, 2)
	for {
		n, err := wr.Read(buf)
		if err != io.EOF {
			if !assert.NoError(t, err) {
				t.FailNow()
			}
		}
		if n == 0 || err == io.EOF {
			break
		}
	}

	log, _ := zap.NewDevelopment()
	log.Info("Test", zap.Array("iolog", iozap.IOLog(wr.Log())))
	time.Sleep(100*time.Millisecond)
	log.Sync()
}
