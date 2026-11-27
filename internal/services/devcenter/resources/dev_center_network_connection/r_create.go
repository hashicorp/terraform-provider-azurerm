// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_network_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	devcenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := devcenterClient.From(metadata.Client).V20250201.NetworkConnections
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var model DevCenterNetworkConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := networkconnections.NewNetworkConnectionID(subscriptionId, model.ResourceGroupName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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

			parameters := networkconnections.NetworkConnection{
				Location: location.Normalize(model.Location),
				Properties: &networkconnections.NetworkProperties{
					DomainJoinType: networkconnections.DomainJoinType(model.DomainJoinType),
					SubnetId:       pointer.To(model.SubnetId),
				},
				Tags: pointer.To(model.Tags),
			}

			if v := model.DomainName; v != "" {
				parameters.Properties.DomainName = pointer.To(v)
			}

			if v := model.DomainPassword; v != "" {
				parameters.Properties.DomainPassword = pointer.To(v)
			}

			if v := model.DomainUsername; v != "" {
				parameters.Properties.DomainUsername = pointer.To(v)
			}

			if v := model.OrganizationUnit; v != "" {
				parameters.Properties.OrganizationUnit = pointer.To(v)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
