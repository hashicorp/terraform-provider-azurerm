// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_managed_private_endpoint

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryManagedPrivateEndpointCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.ManagedPrivateEndpoints
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	managedVirtualNetworkName, err := helpers.GetManagedVirtualNetworkName(ctx, managedVirtualNetworksClient, dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName)
	if err != nil {
		return err
	}
	if managedVirtualNetworkName == nil {
		return fmt.Errorf("managed Private endpoints are only available after managed virtual network for %s is enabled", dataFactoryId)
	}

	id := managedprivateendpoints.NewManagedPrivateEndpointID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, *managedVirtualNetworkName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := getManagedPrivateEndpoint(ctx, client, id.SubscriptionId, id.ResourceGroupName, id.FactoryName, id.ManagedVirtualNetworkName, id.ManagedPrivateEndpointName)
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
		if existing != nil {
			return tf.ImportAsExistsError("azurerm_data_factory_managed_private_endpoint", id.ID())
		}
	}

	targetResourceId := d.Get("target_resource_id").(string)
	subResourceName := d.Get("subresource_name").(string)
	fqdns := d.Get("fqdns").([]any)

	if _, err := privatelinkservices.ParsePrivateLinkServiceID(targetResourceId); err == nil {
		if len(subResourceName) > 0 {
			return fmt.Errorf("`subresource_name` should not be specified when target resource is `Private Link Service`")
		}

		if len(fqdns) == 0 {
			return fmt.Errorf("`fqdns` should be specified when target resource is `Private Link Service`")
		}
	} else {
		if len(strings.TrimSpace(subResourceName)) < 3 {
			return fmt.Errorf("`subresource_name` must be at least 3 characters in length")
		}

		if len(fqdns) > 0 {
			return fmt.Errorf("`fqdns` should not be specified for the target resource: %q", targetResourceId)
		}
	}

	payload := managedprivateendpoints.ManagedPrivateEndpointResource{
		Properties: managedprivateendpoints.ManagedPrivateEndpoint{
			PrivateLinkResourceId: pointer.To(targetResourceId),
		},
	}

	if len(subResourceName) > 0 {
		payload.Properties.GroupId = pointer.To(subResourceName)
	}

	if len(fqdns) > 0 {
		payload.Properties.Fqdns = pluginsdk.ExpandStringSlice(fqdns)
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload, managedprivateendpoints.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Provisioning"},
		Target:     []string{"Succeeded"},
		Refresh:    getManagedPrivateEndpointProvisionStatus(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be created: %+v", id.ID(), err)
	}

	return resourceDataFactoryManagedPrivateEndpointRead(d, meta)
}
