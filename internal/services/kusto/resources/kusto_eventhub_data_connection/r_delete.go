// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_eventhub_data_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoEventHubDataConnectionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnections.ParseDataConnectionID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.DataConnectionName, id.ResourceGroupName, id.ClusterName, id.DatabaseName, err)
	}

	return nil
}
