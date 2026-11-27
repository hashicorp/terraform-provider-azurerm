// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVmwareClusterUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Vmware.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	clusterUpdate := clusters.ClusterUpdate{
		Properties: &clusters.ClusterUpdateProperties{},
	}
	if d.HasChange("cluster_node_count") {
		clusterUpdate.Properties.ClusterSize = pointer.To(int64(d.Get("cluster_node_count").(int)))
	}

	if err := client.UpdateThenPoll(ctx, *id, clusterUpdate); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	return resourceVmwareClusterRead(d, meta)
}
