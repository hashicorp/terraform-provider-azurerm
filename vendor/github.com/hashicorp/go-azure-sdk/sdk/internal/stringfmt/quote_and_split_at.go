// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringfmt

import (
	"fmt"
	"math"
)

// QuoteAndSplitString cuts a string to the specified number of characters and then quotes it
// to enable it to be output as an semi-formatted error message.
func QuoteAndSplitString(message, quoteChar string, characters int) []string {
	lines := make([]string, 0)

	numberOfLines := int(math.Ceil(float64(len(message) / (characters * 1.0))))
	for i := 0; i <= numberOfLines; i++ {
		start := characters * i
		remainingString := message[start:]
		end := characters
		if len(remainingString) < characters {
			end = len(remainingString)
		}
		line := remainingString[0:end]
		lines = append(lines, fmt.Sprintf("%s %s", quoteChar, line))
	}

	return lines
}
