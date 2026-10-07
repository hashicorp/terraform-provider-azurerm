// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ctyhelpers

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
)

// ConstructCtyPath takes a string and converts it to a `cty.Path` for use with `GetRawConfigAt`
// e.g. `identity.0.type`
func ConstructCtyPath(key string) cty.Path {
	p := cty.Path{}

	for segment := range strings.SplitSeq(key, ".") {
		if n, err := strconv.Atoi(segment); err == nil {
			p = p.IndexInt(n)
			continue
		}
		p = p.GetAttr(segment)
	}

	return p
}
