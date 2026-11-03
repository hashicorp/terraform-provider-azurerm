// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package function_app_function

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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
				var confJSON any
				if err = json.Unmarshal([]byte(appFunction.ConfigJSON), &confJSON); err != nil {
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
				Refresh: func() (result any, state string, err error) {
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
