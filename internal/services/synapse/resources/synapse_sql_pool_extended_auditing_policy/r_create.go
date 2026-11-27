// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_extended_auditing_policy

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSqlPoolExtendedAuditingPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.SqlPoolExtendedBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sqlPoolId, err := parse.SqlPoolID(d.Get("sql_pool_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSqlPoolExtendedAuditingPolicyID(sqlPoolId.SubscriptionId, sqlPoolId.ResourceGroup, sqlPoolId.WorkspaceName, sqlPoolId.Name, "default")

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of %s: %+v", id, err)
				}
			}

			// if state is not disabled, we should flag to import it.
			if !response.WasNotFound(existing.Response.Response) {
				if props := existing.ExtendedSQLPoolBlobAuditingPolicyProperties; props != nil && props.State != synapse.BlobAuditingPolicyStateDisabled {
					return tf.ImportAsExistsError("azurerm_synapse_sql_pool_extended_auditing_policy", id.ID())
				}
			}
		}
	}

	params := synapse.ExtendedSQLPoolBlobAuditingPolicy{
		ExtendedSQLPoolBlobAuditingPolicyProperties: &synapse.ExtendedSQLPoolBlobAuditingPolicyProperties{
			State:                       synapse.BlobAuditingPolicyStateEnabled,
			StorageEndpoint:             pointer.To(d.Get("storage_endpoint").(string)),
			IsStorageSecondaryKeyInUse:  pointer.To(d.Get("storage_account_access_key_is_secondary").(bool)),
			RetentionDays:               pointer.To(int32(d.Get("retention_in_days").(int))),
			IsAzureMonitorTargetEnabled: pointer.To(d.Get("log_monitoring_enabled").(bool)),
		},
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		params.StorageAccountAccessKey = pointer.To(v.(string))
	}

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceSynapseSqlPoolExtendedAuditingPolicyRead(d, meta)
}
