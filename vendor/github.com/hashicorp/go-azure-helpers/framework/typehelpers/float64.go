package typehelpers

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type WrappedFloat64Default struct {
	Desc     *string
	Markdown *string
	Value    float64
}

var _ defaults.Float64 = WrappedFloat64Default{}

// NewWrappedFloat64Default is a helper function to return a new defaults.Float64 implementation for any type that
// implements the Go string type.
func NewWrappedFloat64Default[T ~float64](value T) WrappedFloat64Default {
	return WrappedFloat64Default{
		Value: float64(value),
	}
}

func (w WrappedFloat64Default) Description(_ context.Context) string {
	return pointer.From(w.Desc)
}

func (w WrappedFloat64Default) MarkdownDescription(_ context.Context) string {
	return pointer.From(w.Markdown)
}

func (w WrappedFloat64Default) DefaultFloat64(_ context.Context, _ defaults.Float64Request, response *defaults.Float64Response) {
	response.PlanValue = basetypes.NewFloat64Value(w.Value)
}
