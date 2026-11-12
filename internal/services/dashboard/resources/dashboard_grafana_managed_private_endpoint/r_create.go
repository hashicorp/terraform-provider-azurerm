// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana_managed_private_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedprivateendpointmodels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedPrivateEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Dashboard.ManagedPrivateEndpointsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			grafanaId, err := managedgrafanas.ParseGrafanaID(model.GrafanaId)
			if err != nil {
				return err
			}
			id := managedprivateendpointmodels.NewManagedPrivateEndpointID(subscriptionId, grafanaId.ResourceGroupName, grafanaId.GrafanaName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ManagedPrivateEndpointsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			props := managedprivateendpointmodels.ManagedPrivateEndpointModel{
				Location: location.Normalize(model.Location),
				Name:     &model.Name,
				Properties: &managedprivateendpointmodels.ManagedPrivateEndpointModelProperties{
					GroupIds:                  &model.GroupIds,
					PrivateLinkResourceId:     &model.PrivateLinkResourceId,
					PrivateLinkResourceRegion: &model.PrivateLinkResourceRegion,
					RequestMessage:            &model.RequestMessage,
					PrivateLinkServiceURL:     &model.PrivateLinkServiceURL,
				},
				Tags: &model.Tags,
			}

			if err := client.ManagedPrivateEndpointsCreateCallbackThenPoll(ctx, id, props, metadata.SetIDAndIdentityCallback(&id)); err != nil {
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
