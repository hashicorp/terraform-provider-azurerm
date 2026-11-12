// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DashboardGrafanaModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Dashboard.GrafanaResourceClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := managedgrafanas.NewGrafanaID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GrafanaGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := expandCreateForDashboardGrafana(model)

			properties.Identity = expandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))

			if err := client.GrafanaCreateCallbackThenPoll(ctx, id, *properties, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}
