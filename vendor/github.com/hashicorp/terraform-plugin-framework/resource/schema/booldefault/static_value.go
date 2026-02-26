// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package booldefault

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StaticBool returns a static boolean value default handler.
//
// Use StaticBool if a static default value for a boolean should be set.
func StaticBool(defaultVal bool) defaults.Bool {
	return staticBoolDefault{
		defaultVal: defaultVal,
	}
}

// staticBoolDefault is static value default handler that
// sets a value on a boolean attribute.
type staticBoolDefault struct {
	defaultVal bool
}

// Description returns a human-readable description of the default value handler.
func (d staticBoolDefault) Description(_ context.Context) string {
	return fmt.Sprintf("value defaults to %t", d.defaultVal)
}

// MarkdownDescription returns a markdown description of the default value handler.
func (d staticBoolDefault) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value defaults to `%t`", d.defaultVal)
}

// DefaultBool implements the static default value logic.
func (d staticBoolDefault) DefaultBool(_ context.Context, req defaults.BoolRequest, resp *defaults.BoolResponse) {
	resp.PlanValue = types.BoolValue(d.defaultVal)
}
