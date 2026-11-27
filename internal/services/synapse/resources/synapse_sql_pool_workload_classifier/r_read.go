// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_workload_classifier

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSQLPoolWorkloadClassifierRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.SQLPoolWorkloadClassifierClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlPoolWorkloadClassifierID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName, id.WorkloadClassifierName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] synapse %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}
	d.Set("name", id.WorkloadClassifierName)
	d.Set("workload_group_id", parse.NewSqlPoolWorkloadGroupID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName).ID())
	if props := resp.WorkloadClassifierProperties; props != nil {
		d.Set("context", props.Context)
		d.Set("end_time", props.EndTime)
		d.Set("importance", props.Importance)
		d.Set("label", props.Label)
		d.Set("member_name", props.MemberName)
		d.Set("start_time", props.StartTime)
	}
	return nil
}
