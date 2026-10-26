// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateEndpointDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PrivateEndpoints
	dnsZoneGroupsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroups
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateendpoints.ParsePrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	if err := deletePrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsZoneGroupsClient, *id); err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, privateendpoints.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", *id)
	}
	subnetId := ""
	if model := existing.Model; model != nil {
		if props := model.Properties; props != nil {
			if subnet := props.Subnet; subnet != nil && subnet.Id != nil {
				subnetId = *subnet.Id
			}

			// Certain child resources like those for cognitive account lock the parent resource to make sure we don't try and update the parent when it's not ready.
			// Due to that, we'll lock on that resource id to try and prevent those type of errors.
			privateLinkServiceConnectionIds := make([]string, 0)
			if privateLinkServiceConnections := props.PrivateLinkServiceConnections; privateLinkServiceConnections != nil {
				for _, connection := range *privateLinkServiceConnections {
					if connectionProps := connection.Properties; connectionProps != nil {
						privateLinkServiceConnectionIds = append(privateLinkServiceConnectionIds, pointer.From(connectionProps.PrivateLinkServiceId))
					}
				}
			}
			if manualPrivateLinkServiceConnections := props.ManualPrivateLinkServiceConnections; manualPrivateLinkServiceConnections != nil {
				for _, connection := range *manualPrivateLinkServiceConnections {
					if connectionProps := connection.Properties; connectionProps != nil {
						privateLinkServiceConnectionIds = append(privateLinkServiceConnectionIds, pointer.From(connectionProps.PrivateLinkServiceId))
					}
				}
			}

			locks.MultipleByID(pointer.To(privateLinkServiceConnectionIds))
			defer locks.UnlockMultipleByID(pointer.To(privateLinkServiceConnectionIds))
		}
	}
	if subnetId == "" {
		// this also captures `model.Properties` being nil below, since otherwise we wouldn't get the Subnet
		return fmt.Errorf("retrieving existing %s: `model.Properties.Subnet.Id` was nil", *id)
	}

	cosmosDbResIds := getCosmosDbResIdInPrivateServiceConnections(existing.Model.Properties)
	for _, cosmosDbResId := range cosmosDbResIds {
		locks.ByName(cosmosDbResId, "azurerm_private_endpoint")
		//goland:noinspection GoDeferInLoop
		defer locks.UnlockByName(cosmosDbResId, "azurerm_private_endpoint")
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
