// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package typehelpers

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type WrappedInt64Default struct {
	Desc     *string
	Markdown *string
	Value    int64
}

var _ defaults.Int64 = WrappedInt64Default{}

func NewWrappedInt64Default[T ~int](value T) WrappedInt64Default {
	return WrappedInt64Default{
		Value: int64(value),
	}
}

func (w WrappedInt64Default) Description(_ context.Context) string {
	return pointer.From(w.Desc)
}

func (w WrappedInt64Default) MarkdownDescription(_ context.Context) string {
	return pointer.From(w.Markdown)
}

func (w WrappedInt64Default) DefaultInt64(ctx context.Context, request defaults.Int64Request, response *defaults.Int64Response) {
	d := basetypes.NewInt64Value(w.Value)
	response.PlanValue = d
}
