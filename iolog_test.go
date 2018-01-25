package iolog_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/albenik/iolog"
	"github.com/stretchr/testify/assert"
)

var testdata = bytes.Repeat([]byte{'A'}, 8)

func TestWrapper_Read(t *testing.T) {
	r := bytes.NewReader(testdata)
	wr := iolog.WrapReader(r)
	dst := make([]byte, len(testdata))
	n, err := wr.Read(dst)
	assert.NoError(t, err)
	assert.Equal(t, len(testdata), n)
	assert.Equal(t, testdata, dst)
}

func TestWrapper_Write(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, len(testdata)))
	wr := iolog.WrapWriter(buf)
	n, err := wr.Write(testdata)
	assert.NoError(t, err)
	assert.Equal(t, len(testdata), n)
	assert.Equal(t, testdata, buf.Bytes())
}

func TestWrapper_Close(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	wr := iolog.WrapReadCloser(ioutil.NopCloser(buf))
	err := wr.Close()
	assert.NoError(t, err)
}
