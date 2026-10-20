// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_managed_private_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/managedprivateendpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterManagedPrivateEndpointCreateUpdate(d *schema.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ClusterManagedPrivateEndpointClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := managedprivateendpoints.NewManagedPrivateEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			managedPrivateEndpoint, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(managedPrivateEndpoint.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(managedPrivateEndpoint.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_kusto_cluster_managed_private_endpoint", id.ID())
			}
		}
	}

	managedPrivateEndpoint := managedprivateendpoints.ManagedPrivateEndpoint{
		Properties: &managedprivateendpoints.ManagedPrivateEndpointProperties{
			PrivateLinkResourceId: d.Get("private_link_resource_id").(string),
			GroupId:               d.Get("group_id").(string),
		},
	}

	if v, ok := d.GetOk("private_link_resource_region"); ok {
		managedPrivateEndpoint.Properties.PrivateLinkResourceRegion = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("request_message"); ok {
		managedPrivateEndpoint.Properties.RequestMessage = pointer.To(v.(string))
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, managedPrivateEndpoint, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, managedPrivateEndpoint); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceKustoClusterManagedPrivateEndpointRead(d, meta)
}
