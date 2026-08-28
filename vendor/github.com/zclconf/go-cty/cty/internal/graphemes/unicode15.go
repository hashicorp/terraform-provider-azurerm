//go:build !go1.27
// +build !go1.27

package graphemes

import (
	v15 "github.com/apparentlymart/go-textseg/v15/textseg"
)

func ScanGraphemeClusters(data []byte, atEOF bool) (int, []byte, error) {
	return v15.ScanGraphemeClusters(data, atEOF)
}
