// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_managed_private_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedPrivateEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := streamanalyticsClient.From(metadata.Client).EndpointsClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			id := privateendpoints.NewPrivateEndpointID(subscriptionId, model.ResourceGroup, model.StreamAnalyticsCluster, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			props := privateendpoints.PrivateEndpoint{
				Properties: &privateendpoints.PrivateEndpointProperties{
					ManualPrivateLinkServiceConnections: &[]privateendpoints.PrivateLinkServiceConnection{
						{
							Properties: &privateendpoints.PrivateLinkServiceConnectionProperties{
								PrivateLinkServiceId: pointer.To(model.TargetResourceId),
								GroupIds:             &[]string{model.SubResourceName},
							},
						},
					},
				},
			}

			var opts privateendpoints.CreateOrUpdateOperationOptions
			if _, err := client.CreateOrUpdate(ctx, id, props, opts); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
