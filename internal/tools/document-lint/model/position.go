// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package model

type PosType int

const (
	PosDefault PosType = iota
	PosExample
	PosArgs
	PosAttr
	PosTimeout
	PosImport

	PosOther = 100
)

func (p PosType) String() string {
	return [...]string{
		"Defaul",
		"Example",
		"Args",
		"Attr",
		"Timeout",
		"Import",
	}[p]
}

func (p PosType) IsArgOrAttr() bool {
	return p == PosArgs || p == PosAttr
}
