// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package typehelpers

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type WrappedBoolDefault struct {
	Desc     *string
	Markdown *string
	Value    bool
}

var _ defaults.Bool = WrappedBoolDefault{}

func NewWrappedBoolDefault[T ~bool](value T) WrappedBoolDefault {
	return WrappedBoolDefault{
		Value: bool(value),
	}
}

func (w WrappedBoolDefault) Description(_ context.Context) string {
	return pointer.From(w.Desc)
}

func (w WrappedBoolDefault) MarkdownDescription(_ context.Context) string {
	return pointer.From(w.Markdown)
}

func (w WrappedBoolDefault) DefaultBool(_ context.Context, _ defaults.BoolRequest, response *defaults.BoolResponse) {
	d := basetypes.NewBoolValue(w.Value)
	response.PlanValue = d
}
