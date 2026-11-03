// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_managed_private_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/managedprivateendpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoClusterManagedPrivateEndpointDelete(d *schema.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ClusterManagedPrivateEndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
