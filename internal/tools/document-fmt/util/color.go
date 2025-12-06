// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import "github.com/fatih/color"

var (
	Green      = color.New(color.FgGreen).Sprint
	GreenBold  = color.New(color.FgGreen, color.Bold).Sprint
	Red        = color.New(color.FgRed).Sprint
	RedBold    = color.New(color.FgRed, color.Bold).Sprint
	Yellow     = color.New(color.FgYellow).Sprint
	YellowBold = color.New(color.FgYellow, color.Bold).Sprint
	Bold       = color.New(color.Bold).Sprint
)
