// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import "io"

type usefulWriter interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
}

func adapt(w io.Writer) usefulWriter {
	if uw, is := w.(usefulWriter); is {
		return uw
	}
	return &uw{w: w}
}

type uw struct {
	w   io.Writer
	sum int
	err error
}

func (u *uw) Write(p []byte) (n int, err error) {
	if u.err != nil {
		return 0, u.err
	}
	n, err = u.w.Write(p)
	u.sum += n
	u.err = err
	return n, err
}

func (u *uw) WriteString(s string) (n int, err error) {
	return u.Write([]byte(s))
}

func (u *uw) WriteByte(b byte) error {
	_, err := u.Write([]byte{b})
	return err
}

func uwSum(u usefulWriter) (int64, error) {
	if buf, ok := u.(*uw); ok {
		return int64(buf.sum), buf.err
	}
	return 0, nil
}
