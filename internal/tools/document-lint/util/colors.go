// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import "github.com/fatih/color"

var (
	Bold       = color.New(color.Bold).Sprint
	ItalicCode = color.New(color.Italic, color.FgCyan).Sprint
	FormatCode = color.New(color.FgMagenta).Sprint
	Blue       = color.New(color.FgBlue).Sprint
	IssueLine  = color.New(color.FgYellow).Sprint
	FixedCode  = color.New(color.FgGreen).Sprint
)
