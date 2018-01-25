# iolog

[![Build Status](https://travis-ci.org/albenik/iolog.svg?branch=master)](https://travis-ci.org/albenik/iolog)

Golang io.ReadWriteCloser wrapper with log buffer

Used for logging read/write operation (with serial devices in my case) which to noise for direct logging but very helpful in case of hardware error.
