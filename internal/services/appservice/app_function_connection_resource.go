package appservice

import (
	"context"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AppFunctionConnectionResource struct{}

type AppFunctionConnectionModel struct {
	AppID               string `tfschema:"app_id"`
	LinkedFunctionAppID string `tfschema:"linked_function_app_id"`
	Enabled             bool   `tfschema:"enabled"`
}

var _ sdk.ResourceWithUpdate = AppFunctionConnectionResource{}

func (r AppFunctionConnectionResource) ModelObject() interface{} {
	return &AppFunctionConnectionModel{}
}

func (r AppFunctionConnectionResource) ResourceType() string {
	return "azurerm_app_function_connection"
}

func (r AppFunctionConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AppFunctionID
}

func (r AppFunctionConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppID,
			Description:  "The Function App ID to which to attach the function.",
		},

		"linked_function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppID,
			Description:  "The ID of the Function App to link to the Function App referenced by `function_app_id`.",
		},
	}
}

func (r AppFunctionConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AppFunctionConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Create Func
			// TODO - Don't forget to set the ID! e.g. metadata.SetID(id)
			return nil
		},
	}
}

func (r AppFunctionConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Read Func
			return nil
		},
	}
}

func (r AppFunctionConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}

func (r AppFunctionConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Update Func
			return nil
		},
	}
}
