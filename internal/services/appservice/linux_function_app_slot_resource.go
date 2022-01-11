package appservice

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinuxFunctionAppSlotResource struct{}

type LinuxFunctionAppSlotModel struct {
	// TODO - Schema in Go Types format here
}

var _ sdk.ResourceWithUpdate = LinuxFunctionAppSlotResource{}

func (r LinuxFunctionAppSlotResource) ModelObject() interface{} {
	return &LinuxFunctionAppSlotModel{}
}

func (r LinuxFunctionAppSlotResource) ResourceType() string {
	return "azurerm_linux_function_app_slot"
}

func (r LinuxFunctionAppSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	panic("Implement me") // TODO - Add Validation func return here
}

func (r LinuxFunctionAppSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This sections is for configurable items, `Required: true` items first, followed by `Optional: true`,
			both in alphabetical order
		*/
	}
}

func (r LinuxFunctionAppSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
			datasource that can be used as outputs or passed programmatically to other resources or data sources.
		*/
	}
}

func (r LinuxFunctionAppSlotResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Create Func
			// TODO - Don't forget to set the ID! e.g. metadata.SetID(id)
			return nil
		},
	}
}

func (r LinuxFunctionAppSlotResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Read Func
			return nil
		},
	}
}

func (r LinuxFunctionAppSlotResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}

func (r LinuxFunctionAppSlotResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Update Func
			return nil
		},
	}
}
