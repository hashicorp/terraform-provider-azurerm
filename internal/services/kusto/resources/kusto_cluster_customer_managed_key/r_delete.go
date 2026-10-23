// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterID, err := commonids.ParseKustoClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(clusterID.KustoClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(clusterID.KustoClusterName, "azurerm_kusto_cluster")

	// confirm it still exists prior to trying to update it, else we'll get an error
	cluster, err := client.Get(ctx, *clusterID)
	if err != nil {
		if response.WasNotFound(cluster.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *clusterID, err)
	}

	// Since this isn't a real object, just modifying an existing object
	// "Delete" doesn't really make sense it should really be a "Revert to Default"
	// So instead of the Delete func actually deleting the Kusto Cluster I am
	// making it reset the Kusto cluster to its default state
	props := clusters.ClusterUpdate{
		Properties: &clusters.ClusterProperties{
			KeyVaultProperties: &clusters.KeyVaultProperties{},
		},
	}

	if err = client.UpdateThenPoll(ctx, *clusterID, props, clusters.DefaultUpdateOperationOptions()); err != nil {
		return fmt.Errorf("removing Customer Managed Key for %s: %+v", clusterID, err)
	}

	return nil
}
