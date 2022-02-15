package appservice

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AppFunctionResource struct{}

type AppFunctionModel struct {
	// TODO - Schema in Go Types format here
}

var _ sdk.ResourceWithUpdate = AppFunctionResource{}

func (r AppFunctionResource) ModelObject() interface{} {
	return &AppFunctionModel{}
}

func (r AppFunctionResource) ResourceType() string {
	return "azurerm_app_function"
}

func (r AppFunctionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	panic("Implement me") // TODO - Add Validation func return here
}

func (r AppFunctionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This sections is for configurable items, `Required: true` items first, followed by `Optional: true`,
			both in alphabetical order
		*/
	}
}

func (r AppFunctionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
			datasource that can be used as outputs or passed programmatically to other resources or data sources.
		*/
	}
}

func (r AppFunctionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Create Func
			// TODO - Don't forget to set the ID! e.g. metadata.SetID(id)
			return nil
		},
	}
}

func (r AppFunctionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Read Func
			return nil
		},
	}
}

func (r AppFunctionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}

func (r AppFunctionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Update Func
			return nil
		},
	}
}
