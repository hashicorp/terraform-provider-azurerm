// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_extended_auditing_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/blobauditing"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerExtendedAuditingPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.BlobAuditingPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.ExtendedServerBlobAuditingPoliciesGet(ctx, *serverId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("retrieving MsSql Server Extended Auditing Policy %s: %+v", serverId, err)
				}
			}

			// if state is not disabled, we should import it.
			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" && existing.Model.Properties != nil && existing.Model.Properties.State != blobauditing.BlobAuditingPolicyStateDisabled {
				return tf.ImportAsExistsError("azurerm_mssql_server_extended_auditing_policy", *existing.Model.Id)
			}
		}
	}

	params := blobauditing.ExtendedServerBlobAuditingPolicy{
		Properties: &blobauditing.ExtendedServerBlobAuditingPolicyProperties{
			StorageEndpoint:             pointer.To(d.Get("blob_storage_endpoint").(string)),
			IsStorageSecondaryKeyInUse:  pointer.To(d.Get("storage_account_access_key_is_secondary").(bool)),
			RetentionDays:               pointer.To(int64(d.Get("retention_in_days").(int))),
			IsAzureMonitorTargetEnabled: pointer.To(d.Get("log_monitoring_enabled").(bool)),
		},
	}

	if d.Get("enabled").(bool) {
		params.Properties.State = blobauditing.BlobAuditingPolicyStateEnabled
	} else {
		params.Properties.State = blobauditing.BlobAuditingPolicyStateDisabled
	}

	if v, ok := d.GetOk("storage_account_subscription_id"); ok {
		params.Properties.StorageAccountSubscriptionId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		params.Properties.StorageAccountAccessKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("predicate_expression"); ok {
		params.Properties.PredicateExpression = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("audit_actions_and_groups"); ok && len(v.([]any)) > 0 {
		params.Properties.AuditActionsAndGroups = pluginsdk.ExpandStringSlice(v.([]any))
	}

	if err = client.ExtendedServerBlobAuditingPoliciesCreateOrUpdateThenPoll(ctx, *serverId, params); err != nil {
		return fmt.Errorf("creating MsSql Server Extended Auditing Policy %s: %+v", serverId, err)
	}

	id := parse.NewServerExtendedAuditingPolicyID(subscriptionId, serverId.ResourceGroupName, serverId.ServerName, "default")

	d.SetId(id.ID())

	return resourceMsSqlServerExtendedAuditingPolicyRead(d, meta)
}
