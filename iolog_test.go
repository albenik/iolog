package iolog_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/albenik/iolog"
	"github.com/stretchr/testify/assert"
)

var testdata = []byte{1, 2, 3, 4, 5, 6, 7, 8}

func TestWrapper_Read(t *testing.T) {
	r := bytes.NewReader(testdata)
	wr := iolog.WrapReader(r)
	dst := make([]byte, len(testdata))
	n, err := wr.Read(dst)

	assert.NoError(t, err)
	assert.Equal(t, len(testdata), n)
	assert.Equal(t, testdata, dst)
	if !assert.Len(t, wr.Log(), 1) {
		t.FailNow()
	}
	assert.True(t, strings.HasPrefix(wr.Log()[0].String(), "read [01 02 03 04 05 06 07 08]"))
}

func TestWrapper_Write(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, len(testdata)))
	wr := iolog.WrapWriter(buf)
	n, err := wr.Write(testdata)
	assert.NoError(t, err)
	assert.Equal(t, len(testdata), n)
	assert.Equal(t, testdata, buf.Bytes())
	if !assert.Len(t, wr.Log(), 1) {
		t.FailNow()
	}
	assert.True(t, strings.HasPrefix(wr.Log()[0].String(), "write [01 02 03 04 05 06 07 08]"))
}

func TestWrapper_Close(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	wr := iolog.WrapReadCloser(ioutil.NopCloser(buf))
	err := wr.Close()
	assert.NoError(t, err)
	if !assert.Len(t, wr.Log(), 1) {
		t.FailNow()
	}
	assert.True(t, strings.HasPrefix(wr.Log()[0].String(), "close []"))
}

func TestWrapper_AppendCustomLogRecord(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	wr := iolog.WrapReadCloser(ioutil.NopCloser(buf))

	wr.AppendLogRecord(&iolog.Record{Operation: "custom_op", Data: []byte{7, 7, 7}})

	if !assert.Len(t, wr.Log(), 1) {
		t.FailNow()
	}

	assert.True(t, strings.HasPrefix(wr.Log()[0].String(), "custom_op [07 07 07] ("))
	assert.True(t, strings.HasSuffix(wr.Log()[0].String(), "error: <nil>"))

	wr.LastLogRecord().Interface = "test string"
	wr.LastLogRecord().Error = errors.New("test error")
	assert.True(t, strings.HasPrefix(wr.Log()[0].String(), "custom_op [07 07 07] \"test string\" ("))
	assert.True(t, strings.HasSuffix(wr.Log()[0].String(), "error: test error"))
}
