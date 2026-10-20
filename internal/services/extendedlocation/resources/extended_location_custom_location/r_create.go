// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package extended_location_custom_location

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ExtendedLocationCustomLocationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.ExtendedLocation.CustomLocationsClient

			id := customlocations.NewCustomLocationID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			customLocationProps := customlocations.CustomLocationProperties{
				ClusterExtensionIds: pointer.To(model.ClusterExtensionIds),
				DisplayName:         pointer.To(model.DisplayName),
				HostResourceId:      pointer.To(model.HostResourceId),
				HostType:            pointer.ToEnum[customlocations.HostType](model.HostType),
				Namespace:           pointer.To(model.Namespace),
			}

			if len(model.Authentication) > 0 {
				auth := model.Authentication[0]
				customLocationProps.Authentication = &customlocations.CustomLocationPropertiesAuthentication{
					Type:  pointer.To(auth.Type),
					Value: pointer.To(auth.Value),
				}
			}

			props := customlocations.CustomLocation{
				Id:         pointer.To(id.ID()),
				Location:   location.Normalize(model.Location),
				Name:       pointer.To(model.Name),
				Properties: pointer.To(customLocationProps),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, props, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
