// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_extended_auditing_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceExtendedAuditingPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).WorkspaceExtendedBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.WorkspaceExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] synapse %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("synapse_workspace_id", workspaceId.ID())

	if props := resp.ExtendedServerBlobAuditingPolicyProperties; props != nil {
		d.Set("storage_endpoint", props.StorageEndpoint)
		d.Set("storage_account_access_key_is_secondary", props.IsStorageSecondaryKeyInUse)
		d.Set("retention_in_days", props.RetentionDays)
		d.Set("log_monitoring_enabled", props.IsAzureMonitorTargetEnabled)
	}

	return nil
}
