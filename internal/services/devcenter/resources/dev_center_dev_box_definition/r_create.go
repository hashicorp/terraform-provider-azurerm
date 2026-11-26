// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_dev_box_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevBoxDefinitions
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterDevBoxDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := devboxdefinitions.ParseDevCenterID(model.DevCenterId)
			if err != nil {
				return err
			}

			id := devboxdefinitions.NewDevCenterDevBoxDefinitionID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			hs := devboxdefinitions.HibernateSupportDisabled
			if model.HibernateSupport {
				hs = devboxdefinitions.HibernateSupportEnabled
			}

			parameters := devboxdefinitions.DevBoxDefinition{
				Location: location.Normalize(model.Location),
				Properties: &devboxdefinitions.DevBoxDefinitionProperties{
					ImageReference: &devboxdefinitions.ImageReference{
						Id: pointer.To(model.ImageReferenceId),
					},
					HibernateSupport: pointer.To(hs),
					Sku:              expandDevCenterDevBoxDefinitionSku(model.SkuName),
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandDevCenterDevBoxDefinitionSku(input string) *devboxdefinitions.Sku {
	if input == "" {
		return nil
	}

	return &devboxdefinitions.Sku{
		Name: input,
	}
}
