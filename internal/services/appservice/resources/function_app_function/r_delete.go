// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package function_app_function

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseFunctionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			appId := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName)

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			deleteWait := &pluginsdk.StateChangeConf{
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
