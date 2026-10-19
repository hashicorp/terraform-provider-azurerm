// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_resource_bridge_appliance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApplianceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ArcResourceBridge.AppliancesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := appliances.NewApplianceID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			identity, err := identity.ExpandSystemAssignedFromModel(model.Identity)
			if err != nil {
				return err
			}

			parameters := appliances.Appliance{
				Location: location.Normalize(model.Location),
				Properties: &appliances.ApplianceProperties{
					Distro: pointer.To(model.Distro),
					InfrastructureConfig: &appliances.AppliancePropertiesInfrastructureConfig{
						Provider: pointer.To(model.Provider),
					},
				},
				Tags: tags.Expand(model.Tags),
			}

			parameters.Identity = identity

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			// since the public key could not be set during creation, update after creation
			if model.PublicKeyBase64 != "" {
				parameters.Properties.PublicKey = pointer.To(model.PublicKeyBase64)

				if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
