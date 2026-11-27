// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_shared_private_link_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/sharedprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	searchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SharedPrivateLinkServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := searchClient.From(metadata.Client).SearchSharedPrivateLinkResourceClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			searchServiceId, err := services.ParseSearchServiceID(model.SearchServiceId)
			if err != nil {
				return err
			}

			locks.ByID(searchServiceId.ID())
			defer locks.UnlockByID(searchServiceId.ID())

			id := sharedprivatelinkresources.NewSharedPrivateLinkResourceID(subscriptionId, searchServiceId.ResourceGroupName, searchServiceId.SearchServiceName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, sharedprivatelinkresources.GetOperationOptions{})
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing shared private link resource %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := sharedprivatelinkresources.SharedPrivateLinkResource{
				Properties: &sharedprivatelinkresources.SharedPrivateLinkResourceProperties{
					GroupId:               pointer.To(model.SubResourceName),
					PrivateLinkResourceId: pointer.To(model.TargetResourceId),
				},
			}

			if model.RequestMessage != "" {
				parameters.Properties.RequestMessage = pointer.To(model.RequestMessage)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sharedprivatelinkresources.CreateOrUpdateOperationOptions{}, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
		Timeout: 60 * time.Minute,
	}
}
