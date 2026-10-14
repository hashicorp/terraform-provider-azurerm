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
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetFunction(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					if !response.WasBadRequest(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) && !response.WasBadRequest(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var confJSON any
			if err = json.Unmarshal([]byte(appFunction.ConfigJSON), &confJSON); err != nil {
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
				Refresh: func() (result any, state string, err error) {
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

			if err := client.CreateFunctionCallbackThenPoll(ctx, id, fnEnvelope, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
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
