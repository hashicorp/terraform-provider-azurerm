// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_extended_auditing_policy

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceExtendedAuditingPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).WorkspaceExtendedBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewWorkspaceExtendedAuditingPolicyID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, "default")

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of %s: %+v", id, err)
				}
			}

			// if state is not disabled, we should flag it for import
			if !response.WasNotFound(existing.Response.Response) {
				if props := existing.ExtendedServerBlobAuditingPolicyProperties; props != nil && props.State != synapse.BlobAuditingPolicyStateDisabled {
					return tf.ImportAsExistsError("azurerm_synapse_workspace_extended_auditing_policy", id.ID())
				}
			}
		}
	}

	params := synapse.ExtendedServerBlobAuditingPolicy{
		ExtendedServerBlobAuditingPolicyProperties: &synapse.ExtendedServerBlobAuditingPolicyProperties{
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, params)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for creation of %s: %+v", id, err)
		}
	}

	return resourceSynapseWorkspaceExtendedAuditingPolicyRead(d, meta)
}
