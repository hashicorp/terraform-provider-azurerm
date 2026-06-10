// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package util

import "github.com/fatih/color"

var (
	Green     = color.New(color.FgGreen).Sprint
	GreenBold = color.New(color.FgGreen, color.Bold).Sprint
	Red       = color.New(color.FgRed).Sprint
	RedBold   = color.New(color.FgRed, color.Bold).Sprint
)
