// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "testing"

// TestLocalHelperResolution verifies that property values and Elem schemas
// produced by local schema-returning helpers are followed, so block-level rules
// fire on the block and leaf rules see the children with full dotted paths.
func TestLocalHelperResolution(t *testing.T) {
	src := `package x
func r() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"direct":  directBlock(),
			"via_var": viaVarBlock(),
		},
	}
}

func directBlock() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:        pluginsdk.TypeList,
		Optional:    true,
		Description: "d",
		Elem: &pluginsdk.Resource{Schema: func() map[string]*pluginsdk.Schema {
			return map[string]*pluginsdk.Schema{
				"a": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
				"b": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
			}
		}()},
	}
}

func viaVarBlock() *pluginsdk.Schema {
	s := &pluginsdk.Schema{
		Type:        pluginsdk.TypeList,
		Optional:    true,
		Description: "d",
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"c": {Type: pluginsdk.TypeString, Optional: true, Description: "d"},
		}},
	}
	return s
}
`

	// SL006: both helper-composed blocks are all-optional and are flagged at the
	// call site (direct's children come via an IIFE, via_var's via a
	// `s := ...; return s` helper).
	fs6 := runRule(t, sl006, src)
	for _, p := range []string{"direct", "via_var"} {
		if !flagged(fs6, p) {
			t.Errorf("expected SL006 on resolved block %q", p)
		}
	}

	// SL005: resolved children are linted with full dotted paths, and the
	// consumed helper maps are not also linted as standalone roots.
	fs5 := runRule(t, sl005, src)
	for _, p := range []string{"direct.a", "direct.b", "via_var.c"} {
		if !flagged(fs5, p) {
			t.Errorf("expected resolved child %q", p)
		}
	}
	for _, p := range []string{"a", "b", "c"} {
		if flagged(fs5, p) {
			t.Errorf("resolved child should not also appear as a standalone root %q", p)
		}
	}
}
