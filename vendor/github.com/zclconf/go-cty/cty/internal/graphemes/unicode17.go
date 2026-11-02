//go:build go1.27

package graphemes

import (
	v17 "github.com/apparentlymart/go-textseg/v17/textseg"
)

func ScanGraphemeClusters(data []byte, atEOF bool) (int, []byte, error) {
	return v17.ScanGraphemeClusters(data, atEOF)
}
