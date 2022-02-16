package appservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppFunctionResource struct{}

type AppFunctionModel struct {
	Name           string `tfschema:"name"`
	AppID          string `tfschema:"function_app_id"`
	Enabled        bool   `tfschema:"enabled"`
	ConfigJSON     string `tfschema:"config_json"`
	Language       string `tfschema:"language"`
	SecretsFileURL string `tfschema:"secrets_file_url"`
	TestData       string `tfschema:"test_data"`

	ConfigURL         string `tfschema:"config_url"`
	FunctionURL       string `tfschema:"url"`
	InvokeURL         string `tfschema:"invocation_url"`
	ScriptURL         string `tfschema:"script_url"`
	ScriptRootPathURL string `tfschema:"script_root_path_url"`
	TestDataURL       string `tfschema:"test_data_url"`
}

var _ sdk.ResourceWithUpdate = AppFunctionResource{}

func (r AppFunctionResource) ModelObject() interface{} {
	return &AppFunctionModel{}
}

func (r AppFunctionResource) ResourceType() string {
	return "azurerm_app_function"
}

func (r AppFunctionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AppFunctionID
}

func (r AppFunctionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 128), // TODO - proper validation here, "must start with a letter and can contain letters, numbers (0-9), dashes ("-"), and underscores ("_")."
		},

		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppID,
			Description:  "The Function App ID to which this function belongs.",
		},

		"config_json": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsJSON,
			Description:  "The JSON config for this Function.",
		},

		"enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Should this function be enabled. Defaults to `true`",
		},

		"language": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty, // TODO - find the valida list of strings
			Description:  "The language the Function is written in.",
		},

		"test_data": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty, // Can an empty string be valid here?
			Description:  "The test data for the function.",
		},
	}
}

func (r AppFunctionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"config_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The URL of the configuration JSON.",
		},

		"url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The function URL.",
		},

		"invocation_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The invocation URL.",
		},

		"script_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The script URL.",
		},

		"script_root_path_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Script root path URL.",
		},

		"test_data_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Test data URL.",
		},
		"secrets_file_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The URL for the Secrets File.",
		},
	}
}

func (r AppFunctionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			appFunction := AppFunctionModel{}
			if err := metadata.Decode(&appFunction); err != nil {
				return err
			}

			appId, err := parse.FunctionAppID(appFunction.AppID)
			if err != nil {
				return err
			}

			id := parse.NewAppFunctionID(appId.SubscriptionId, appId.ResourceGroup, appId.SiteName, appFunction.Name)

			existing, err := client.GetFunction(ctx, id.ResourceGroup, id.SiteName, id.FunctionName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var confJSON interface{}
			err = json.Unmarshal([]byte(appFunction.ConfigJSON), &confJSON)
			if err != nil {
				return fmt.Errorf("error preparing config data to send: %+v", err)
			}

			fnEnvelope := web.FunctionEnvelope{
				FunctionEnvelopeProperties: &web.FunctionEnvelopeProperties{
					Config:     confJSON,
					TestData:   utils.String(appFunction.TestData),
					Language:   utils.String(appFunction.Language),
					IsDisabled: utils.Bool(!appFunction.Enabled),
					// Files:      nil, // TODO - Can / should we support this?
				},
			}

			future, err := client.CreateFunction(ctx, id.ResourceGroup, id.SiteName, id.FunctionName, fnEnvelope)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AppFunctionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.AppFunctionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetFunction(ctx, id.ResourceGroup, id.SiteName, id.FunctionName)
			if err != nil || existing.FunctionEnvelopeProperties == nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			appFunc := AppFunctionModel{
				Name:              id.FunctionName,
				AppID:             parse.NewFunctionAppID(id.SubscriptionId, id.ResourceGroup, id.SiteName).ID(),
				ConfigURL:         utils.NormalizeNilableString(existing.ConfigHref),
				Enabled:           !utils.NormaliseNilableBool(existing.IsDisabled),
				FunctionURL:       utils.NormalizeNilableString(existing.Href),
				InvokeURL:         utils.NormalizeNilableString(existing.InvokeURLTemplate),
				Language:          utils.NormalizeNilableString(existing.Language),
				ScriptURL:         utils.NormalizeNilableString(existing.ScriptHref),
				ScriptRootPathURL: utils.NormalizeNilableString(existing.ScriptRootPathHref),
				SecretsFileURL:    utils.NormalizeNilableString(existing.SecretsFileHref),
				TestData:          utils.NormalizeNilableString(existing.TestData),
				TestDataURL:       utils.NormalizeNilableString(existing.TestDataHref),
			}

			metadata.Logger.Infof("[STEBUG] config_json: %#v", existing.Config)
			if configJSON := existing.Config; configJSON != nil {
				raw, err := json.Marshal(configJSON)
				if err != nil {
					metadata.Logger.Infof("[STEBUG] JSON: %+v", err)
				}
				appFunc.ConfigJSON = string(raw)
			}

			return metadata.Encode(&appFunc)
		},
	}
}

func (r AppFunctionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.AppFunctionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if _, err = client.DeleteFunction(ctx, id.ResourceGroup, id.SiteName, id.FunctionName); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

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
