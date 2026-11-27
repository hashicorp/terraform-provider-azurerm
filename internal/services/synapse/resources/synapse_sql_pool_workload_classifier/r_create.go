// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_workload_classifier

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

func resourceSynapseSQLPoolWorkloadClassifierCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).SQLPoolWorkloadClassifierClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	workloadGroupId, err := parse.SqlPoolWorkloadGroupID(d.Get("workload_group_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSqlPoolWorkloadClassifierID(workloadGroupId.SubscriptionId, workloadGroupId.ResourceGroup, workloadGroupId.WorkspaceName, workloadGroupId.SqlPoolName, workloadGroupId.WorkloadGroupName, d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName, id.WorkloadClassifierName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_synapse_sql_pool_workload_classifier", id.ID())
			}
		}
	}

	parameters := synapse.WorkloadClassifier{
		WorkloadClassifierProperties: &synapse.WorkloadClassifierProperties{
			Context:    pointer.To(d.Get("context").(string)),
			EndTime:    pointer.To(d.Get("end_time").(string)),
			Importance: pointer.To(d.Get("importance").(string)),
			Label:      pointer.To(d.Get("label").(string)),
			MemberName: pointer.To(d.Get("member_name").(string)),
			StartTime:  pointer.To(d.Get("start_time").(string)),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName, id.WorkloadClassifierName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	return resourceSynapseSQLPoolWorkloadClassifierRead(d, meta)
}
