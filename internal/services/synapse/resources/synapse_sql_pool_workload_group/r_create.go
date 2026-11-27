// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_workload_group

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSQLPoolWorkloadGroupCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).SQLPoolWorkloadGroupClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	sqlPoolId, err := parse.SqlPoolID(d.Get("sql_pool_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSqlPoolWorkloadGroupID(sqlPoolId.SubscriptionId, sqlPoolId.ResourceGroup, sqlPoolId.WorkspaceName, sqlPoolId.Name, d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_synapse_sql_pool_workload_group", id.ID())
			}
		}
	}

	parameters := synapse.WorkloadGroup{
		WorkloadGroupProperties: &synapse.WorkloadGroupProperties{
			Importance:                   pointer.To(d.Get("importance").(string)),
			MaxResourcePercent:           pointer.To(int32(d.Get("max_resource_percent").(int))),
			MaxResourcePercentPerRequest: pointer.To(d.Get("max_resource_percent_per_request").(float64)),
			MinResourcePercent:           pointer.To(int32(d.Get("min_resource_percent").(int))),
			MinResourcePercentPerRequest: pointer.To(d.Get("min_resource_percent_per_request").(float64)),
		},
	}

	if timeout, ok := d.GetOk("query_execution_timeout_in_seconds"); ok {
		parameters.QueryExecutionTimeout = pointer.To(int32(timeout.(int)))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	return resourceSynapseSQLPoolWorkloadGroupRead(d, meta)
}
