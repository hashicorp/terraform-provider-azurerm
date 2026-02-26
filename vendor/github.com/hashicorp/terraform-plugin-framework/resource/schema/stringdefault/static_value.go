// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringdefault

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StaticString returns a static string value default handler.
//
// Use StaticString if a static default value for a string should be set.
func StaticString(defaultVal string) defaults.String {
	return staticStringDefault{
		defaultVal: defaultVal,
	}
}

// staticStringDefault is static value default handler that
// sets a value on a string attribute.
type staticStringDefault struct {
	defaultVal string
}

// Description returns a human-readable description of the default value handler.
func (d staticStringDefault) Description(_ context.Context) string {
	return fmt.Sprintf("value defaults to %s", d.defaultVal)
}

// MarkdownDescription returns a markdown description of the default value handler.
func (d staticStringDefault) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value defaults to `%s`", d.defaultVal)
}

// DefaultString implements the static default value logic.
func (d staticStringDefault) DefaultString(_ context.Context, req defaults.StringRequest, resp *defaults.StringResponse) {
	resp.PlanValue = types.StringValue(d.defaultVal)
}
