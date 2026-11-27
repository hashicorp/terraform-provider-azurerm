// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_cluster

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	vmwareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwareClusterRead(d *pluginsdk.ResourceData, meta any) error {
	client := vmwareClient.From(meta).ClusterClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("vmware_cloud_id", privateclouds.NewPrivateCloudID(id.SubscriptionId, id.ResourceGroupName, id.PrivateCloudName).ID())

	if model := resp.Model; model != nil {
		d.Set("cluster_node_count", model.Properties.ClusterSize)
		d.Set("cluster_number", model.Properties.ClusterId)
		d.Set("hosts", pluginsdk.FlattenSlice(model.Properties.Hosts))
		d.Set("sku_name", model.Sku.Name)
	}

	return nil
}
