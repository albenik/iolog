package iozap

import "fmt"

type StringerFunc func() string

func (fn StringerFunc) String() string {
	return fn()
}

func BytearrayStringer(bytes []byte) fmt.Stringer {
	return StringerFunc(func() string {
		return fmt.Sprintf("% X", bytes)
	})
}
