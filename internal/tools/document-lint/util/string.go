// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"log"
	"regexp"
	"strings"
)

func FirstCodeValue(line string) string {
	if vals := ExtractCodeValue(line); len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func ExtractCodeValue(line string) (res []string) {
	idx1 := strings.Index(line, "`")
	for idx1 >= 0 {
		idx2 := idx1 + 1 + strings.Index(line[idx1+1:], "`")
		if idx2 > len(line) || idx2 <= idx1 {
			log.Printf("ExtractCodeValue: code mark ` not closed in '%s'", line)
			return
		}
		res = append(res, line[idx1+1:idx2])
		nextIdx := strings.Index(line[idx2+1:], "`")
		if nextIdx < 0 {
			break
		}
		idx1 = idx2 + 1 + nextIdx
	}
	return
}

var timeoutValueReg = regexp.MustCompile(`[0-9]+ (hours?|minutes?)`)

func TimeoutValueIdx(line string) (start, end int) {
	idx := timeoutValueReg.FindStringSubmatchIndex(line)
	if len(idx) >= 2 {
		// search to next ')' if exists
		if idx2 := strings.Index(line[idx[1]:], ")"); idx2 > 0 {
			idx[1] += idx2
		}
		return idx[0], idx[1]
	}
	return
}

func NormalizeResourceName(rt string) string {
	// parse rt as normailize name
	rtBs := []byte(strings.TrimPrefix(rt, "azurerm_"))
	for idx, ch := range rtBs {
		if idx == 0 && (ch >= 'a' && ch <= 'z') {
			rtBs[0] = ch - 32
		}
		if ch == '_' {
			rtBs[idx] = ' '
			if idx+1 < len(rtBs) {
				// convert to uppercase
				if ch2 := rtBs[idx+1]; ch2 >= 'a' && ch2 <= 'z' {
					rtBs[idx+1] = ch2 - 32
				}
			}
		}
	}
	rt = string(rtBs)
	return rt
}

func XPathBase(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '.' {
			return p[i+1:]
		}
	}

	return p
}

func XPathDir(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '.' {
			return p[:i]
		}
	}
	return p
}
