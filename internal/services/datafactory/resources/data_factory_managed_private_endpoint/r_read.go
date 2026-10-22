// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_managed_private_endpoint

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedprivateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryManagedPrivateEndpointRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.ManagedPrivateEndpoints
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, managedprivateendpoints.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ManagedPrivateEndpointName)
	d.Set("data_factory_id", factories.NewFactoryID(subscriptionId, id.ResourceGroupName, id.FactoryName).ID())

	if model := resp.Model; model != nil {
		props := model.Properties
		d.Set("target_resource_id", props.PrivateLinkResourceId)
		d.Set("subresource_name", props.GroupId)
		d.Set("fqdns", pluginsdk.FlattenSlice(props.Fqdns))
	}

	return nil
}

// if ManagedPrivateEndpoint not exist, get rest api will return 400 bad request
// invoke list rets api and then filter by name
func getManagedPrivateEndpoint(ctx context.Context, client *managedprivateendpoints.ManagedPrivateEndpointsClient, subscriptionId, resourceGroupName, factoryName, managedVirtualNetworkName, privateEndpointName string) (*managedprivateendpoints.ManagedPrivateEndpointResource, error) {
	managedVirtualNetworkId := managedprivateendpoints.NewManagedVirtualNetworkID(subscriptionId, resourceGroupName, factoryName, managedVirtualNetworkName)
	iter, err := client.ListByFactoryComplete(ctx, managedVirtualNetworkId)
	if err != nil {
		return nil, err
	}

	for _, item := range iter.Items {
		if item.Name != nil && *item.Name == privateEndpointName {
			return &item, nil
		}
	}
	return nil, nil
}

func getManagedPrivateEndpointProvisionStatus(ctx context.Context, client *managedprivateendpoints.ManagedPrivateEndpointsClient, id managedprivateendpoints.ManagedPrivateEndpointId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		// TODO: it should be possible to remove this function https://github.com/hashicorp/go-azure-sdk/issues/307 has been fixed
		resp, err := client.Get(ctx, id, managedprivateendpoints.DefaultGetOperationOptions())
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		model := resp.Model
		if model == nil {
			return nil, "", fmt.Errorf("retrieving %s: `model` was nil", id)
		}

		if model.Properties.ProvisioningState == nil {
			return nil, "", fmt.Errorf("retrieving %s: `provisioningState` is nil", id)
		}

		return resp, *resp.Model.Properties.ProvisioningState, nil
	}
}
