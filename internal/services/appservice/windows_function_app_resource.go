package appservice 

import (
    "context"
    "time"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WindowsFunctionAppResource struct{}

type WindowsFunctionAppModel struct {
    // TODO - Schema in Go Types format here
}

var _ sdk.ResourceWithUpdate = WindowsFunctionAppResource{}

func (r WindowsFunctionAppResource) ModelObject() interface{} {
    return WindowsFunctionAppModel{}
}

func (r WindowsFunctionAppResource) ResourceType() string {
    return "azurerm_windows_function_app"
}

func (r WindowsFunctionAppResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    panic("Implement me") // TODO - Add Validation func return here
}

func (r WindowsFunctionAppResource) Arguments() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        /*
            TODO - This sections is for configurable items, `Required: true` items first, followed by `Optional: true`,
            both in alphabetical order
        */
    }
}

func (r WindowsFunctionAppResource) Attributes() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        /*
            TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
            datasource that can be used as outputs or passed programmatically to other resources or data sources.
        */
    }
}

func (r WindowsFunctionAppResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // TODO - Create Func
            // TODO - Don't forget to set the ID! e.g. metadata.SetID(id)
            return nil
        },
    }
}

func (r WindowsFunctionAppResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // TODO - Read Func
            return nil
        },
    }
}

func (r WindowsFunctionAppResource) Delete() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // TODO - Delete Func
            return nil
        },
    }
}

func (r WindowsFunctionAppResource) Update() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // TODO - Delete Func
            return nil
        },
    }
}
