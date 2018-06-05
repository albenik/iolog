package iolog_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/albenik/iolog"
)

var testdata = []byte{1, 2, 3, 4, 5, 6, 7, 8}

func TestIOLog_LogIO(t *testing.T) {
	l := iolog.New(128)
	l.Start()
	buf := make([]byte, len(testdata))

	src := bytes.NewReader(testdata)
	n, err := l.LogIO("read", src.Read, buf)

	assert.NoError(t, err)
	assert.Equal(t, len(testdata), n)
	assert.Equal(t, testdata, buf)
	if assert.Equal(t, 1, l.Len()) {
		assert.True(t, strings.HasPrefix(l.LastRecord().String(), "read [01 02 03 04 05 06 07 08]"), l.LastRecord().String())
	}

	dst := bytes.NewBuffer(make([]byte, 8))
	n, err = l.LogIO("write", dst.Write, testdata)
	assert.NoError(t, err)
	assert.Equal(t, len(testdata), n)
	assert.Equal(t, testdata, buf)
	if assert.Equal(t, 2, l.Len()) {
		assert.True(t, strings.HasPrefix(l.LastRecord().String(), "write [01 02 03 04 05 06 07 08]"), l.LastRecord().String())
	}
}

func TestIOLog_LogAny(t *testing.T) {
	l := iolog.New(128)
	l.Start()
	err := l.LogAny("any", func() (interface{}, error) {
		return 777, nil
	})

	assert.NoError(t, err)
	if assert.Equal(t, 1, l.Len()) {
		assert.True(t, strings.HasPrefix(l.LastRecord().String(), "any 777"), "%q", l.LastRecord().String())
	}
}
