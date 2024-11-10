// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FunctionAppFunctionResource struct{}

type FunctionAppFunctionModel struct {
	Name           string          `tfschema:"name"`
	AppID          string          `tfschema:"function_app_id"`
	Enabled        bool            `tfschema:"enabled"`
	ConfigJSON     string          `tfschema:"config_json"`
	Language       string          `tfschema:"language"`
	SecretsFileURL string          `tfschema:"secrets_file_url"`
	TestData       string          `tfschema:"test_data"`
	Files          []FunctionFiles `tfschema:"file"`

	ConfigURL         string `tfschema:"config_url"`
	FunctionURL       string `tfschema:"url"`
	InvokeURL         string `tfschema:"invocation_url"`
	ScriptURL         string `tfschema:"script_url"`
	ScriptRootPathURL string `tfschema:"script_root_path_url"`
	TestDataURL       string `tfschema:"test_data_url"`
}

type FunctionFiles struct {
	Name    string `tfschema:"name"`
	Content string `tfschema:"content"`
}

var _ sdk.ResourceWithUpdate = FunctionAppFunctionResource{}

func (r FunctionAppFunctionResource) ModelObject() interface{} {
	return &FunctionAppFunctionModel{}
}

func (r FunctionAppFunctionResource) ResourceType() string {
	return "azurerm_function_app_function"
}

func (r FunctionAppFunctionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webapps.ValidateFunctionID
}

func (r FunctionAppFunctionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FunctionAppFunctionName,
			Description:  "The name of the function.",
		},

		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
			ForceNew:     true,
			Description:  "The ID of the Function App in which this function should reside.",
		},

		"config_json": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsJSON,
			Description:  "The config for this Function in JSON format.",
		},

		"enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Should this function be enabled. Defaults to `true`.",
		},

		"language": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"CSharp",
				"Custom",
				"Java",
				"Javascript",
				"Python",
				"PowerShell",
				"TypeScript",
			}, false), // TODO - find the valida list of strings
			Description: "The language the Function is written in.",
		},

		"file": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  "The filename of the file to be uploaded.",
					},
					"content": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  "The content of the file.",
					},
				},
			},
		},

		"test_data": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The test data for the function.",
		},
	}
}

func (r FunctionAppFunctionResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r FunctionAppFunctionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			appFunction := FunctionAppFunctionModel{}
			if err := metadata.Decode(&appFunction); err != nil {
				return err
			}

			appId, err := commonids.ParseFunctionAppID(appFunction.AppID)
			if err != nil {
				return err
			}

			id := webapps.NewFunctionID(appId.SubscriptionId, appId.ResourceGroupName, appId.SiteName, appFunction.Name)

			existing, err := client.GetFunction(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				if !response.WasBadRequest(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) && !response.WasBadRequest(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var confJSON interface{}
			err = json.Unmarshal([]byte(appFunction.ConfigJSON), &confJSON)
			if err != nil {
				return fmt.Errorf("error preparing config data to send: %+v", err)
			}

			fnEnvelope := webapps.FunctionEnvelope{
				Properties: &webapps.FunctionEnvelopeProperties{
					Config:     pointer.To(confJSON),
					TestData:   pointer.To(appFunction.TestData),
					Language:   pointer.To(appFunction.Language),
					IsDisabled: pointer.To(!appFunction.Enabled),
					Files:      expandFunctionFiles(appFunction.Files),
				},
			}

			// Check and wait for the Function to have no in flight operations
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			createWait := &pluginsdk.StateChangeConf{
				Pending: []string{"busy", "unknown"},
				Target:  []string{"ready"},
				Refresh: func() (result interface{}, state string, err error) {
					function, err := client.Get(ctx, *appId)
					if err != nil || function.Model == nil || function.Model.Properties == nil {
						return "unknown", "unknown", err
					}
					if function.Model.Properties.InProgressOperationId != nil {
						return "busy", "busy", nil
					}
					return "ready", "ready", nil
				},
				MinTimeout:                30 * time.Second,
				ContinuousTargetOccurence: 2,
				Timeout:                   time.Until(deadline),
			}

			if _, err = createWait.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be ready", *appId)
			}

			locks.ByID(appId.ID())
			defer locks.UnlockByID(appId.ID())

			if err := client.CreateFunctionThenPoll(ctx, id, fnEnvelope); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r FunctionAppFunctionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseFunctionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetFunction(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			appFunc := FunctionAppFunctionModel{
				Name:  id.FunctionName,
				AppID: commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName).ID(),
			}

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					appFunc.ConfigURL = pointer.From(props.ConfigHref)
					appFunc.Enabled = !pointer.From(props.IsDisabled)
					appFunc.FunctionURL = pointer.From(props.Href)
					appFunc.InvokeURL = pointer.From(props.InvokeUrlTemplate)
					appFunc.ScriptURL = pointer.From(props.ScriptHref)
					appFunc.ScriptRootPathURL = pointer.From(props.ScriptRootPathHref)
					appFunc.SecretsFileURL = pointer.From(props.SecretsFileHref)
					appFunc.TestData = pointer.From(props.TestData)
					appFunc.TestDataURL = pointer.From(props.TestDataHref)
					if language, ok := metadata.ResourceData.GetOk("language"); ok {
						appFunc.Language = language.(string)
					}

					if filesRaw, ok := metadata.ResourceData.GetOk("file"); ok {
						files := make([]FunctionFiles, 0)
						for _, v := range filesRaw.([]interface{}) {
							file := v.(map[string]interface{})
							files = append(files, FunctionFiles{
								Name:    file["name"].(string),
								Content: file["content"].(string),
							})
						}
						appFunc.Files = files
					}

					config, err := flattenFunctionFiles(props.Config)
					if err != nil {
						return err
					}
					appFunc.ConfigJSON = pointer.From(config)
				}
			}

			return metadata.Encode(&appFunc)
		},
	}
}

func (r FunctionAppFunctionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseFunctionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			appId := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName)

			metadata.Logger.Infof("deleting %s", *id)

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			deleteWait := &pluginsdk.StateChangeConf{
				Pending: []string{"busy", "unknown"},
				Target:  []string{"ready"},
				Refresh: func() (result interface{}, state string, err error) {
					function, err := client.Get(ctx, appId)
					if err != nil || function.Model == nil || function.Model.Properties == nil || function.Model.Properties.SiteConfig == nil {
						return "unknown", "unknown", err
					}
					if function.Model.Properties.InProgressOperationId != nil {
						return "busy", "busy", nil
					}
					return "ready", "ready", nil
				},
				MinTimeout:                30 * time.Second,
				ContinuousTargetOccurence: 2,
				Timeout:                   time.Until(deadline),
			}

			if _, err = deleteWait.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be settled", *id)
			}

			locks.ByID(appId.ID())
			defer locks.UnlockByID(appId.ID())

			if _, err = client.DeleteFunction(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r FunctionAppFunctionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseFunctionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			appId := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName)

			var appFunction FunctionAppFunctionModel
			if err := metadata.Decode(&appFunction); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.GetFunction(ctx, *id)
			if err != nil || existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("config_json") {
				var confJSON interface{}
				err = json.Unmarshal([]byte(appFunction.ConfigJSON), &confJSON)
				if err != nil {
					return fmt.Errorf("error preparing config data to send: %+v", err)
				}
				model.Properties.Config = pointer.To(confJSON)
			}

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.IsDisabled = pointer.To(!appFunction.Enabled)
			}

			if metadata.ResourceData.HasChange("test_data") {
				model.Properties.TestData = pointer.To(appFunction.TestData)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			updateWait := &pluginsdk.StateChangeConf{
				Pending: []string{"busy", "unknown"},
				Target:  []string{"ready"},
				Refresh: func() (result interface{}, state string, err error) {
					function, err := client.Get(ctx, appId)
					if err != nil || function.Model == nil || function.Model.Properties == nil || function.Model.Properties.SiteConfig == nil {
						return "unknown", "unknown", err
					}
					if function.Model.Properties.InProgressOperationId != nil {
						return "busy", "busy", nil
					}
					return "ready", "ready", nil
				},
				MinTimeout:                30 * time.Second,
				ContinuousTargetOccurence: 2,
				Timeout:                   time.Until(deadline),
			}

			if _, err = updateWait.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be ready", *id)
			}

			locks.ByID(appId.ID())
			defer locks.UnlockByID(appId.ID())

			if err := client.CreateFunctionThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandFunctionFiles(input []FunctionFiles) *map[string]string {
	if input == nil {
		return nil
	}
	result := make(map[string]string)
	for _, v := range input {
		result[v.Name] = v.Content
	}

	return &result
}

func flattenFunctionFiles(input interface{}) (*string, error) {
	if input == nil {
		return nil, nil
	}

	raw, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("could not marshal `config_json`: %+v", err)
	}
	result := string(raw)
	return &result, nil
}
