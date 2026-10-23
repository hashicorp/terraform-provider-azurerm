// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package function_app_function

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
						for _, v := range filesRaw.([]any) {
							file := v.(map[string]any)
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

func flattenFunctionFiles(input any) (*string, error) {
	if input == nil {
		return nil, nil
	}

	raw, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("could not marshal `config_json`: %+v", err)
	}
	return pointer.To(string(raw)), nil
}

type FunctionFiles struct {
	Name    string `tfschema:"name"`
	Content string `tfschema:"content"`
}
