// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_workload_group

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSQLPoolWorkloadGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).SQLPoolWorkloadGroupClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.SqlPoolWorkloadGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] synapse %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}
	d.Set("name", id.WorkloadGroupName)
	d.Set("sql_pool_id", parse.NewSqlPoolID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName).ID())
	if props := resp.WorkloadGroupProperties; props != nil {
		d.Set("importance", props.Importance)
		d.Set("max_resource_percent", props.MaxResourcePercent)
		d.Set("max_resource_percent_per_request", props.MaxResourcePercentPerRequest)
		d.Set("min_resource_percent", props.MinResourcePercent)
		d.Set("min_resource_percent_per_request", props.MinResourcePercentPerRequest)
		d.Set("query_execution_timeout_in_seconds", props.QueryExecutionTimeout)
	}
	return nil
}
