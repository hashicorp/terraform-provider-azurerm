// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_managed_private_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/managedprivateendpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterManagedPrivateEndpointRead(d *schema.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ClusterManagedPrivateEndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ManagedPrivateEndpointName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if resp.Model != nil {
		props := resp.Model.Properties
		if props != nil {
			d.Set("private_link_resource_id", props.PrivateLinkResourceId)
			d.Set("group_id", props.GroupId)

			if props.PrivateLinkResourceRegion != nil {
				d.Set("private_link_resource_region", props.PrivateLinkResourceRegion)
			}

			if props.RequestMessage != nil {
				d.Set("request_message", props.RequestMessage)
			}
		}
	}

	return nil
}
