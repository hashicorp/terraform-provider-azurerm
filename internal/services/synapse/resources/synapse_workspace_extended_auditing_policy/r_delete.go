// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_extended_auditing_policy

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceExtendedAuditingPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.WorkspaceExtendedBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	params := synapse.ExtendedServerBlobAuditingPolicy{
		ExtendedServerBlobAuditingPolicyProperties: &synapse.ExtendedServerBlobAuditingPolicyProperties{
			State: synapse.BlobAuditingPolicyStateDisabled,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, params)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
		}
	}

	return nil
}
