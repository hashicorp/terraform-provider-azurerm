// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_extended_auditing_policy

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseSqlPoolExtendedAuditingPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).SqlPoolExtendedBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.SqlPoolExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	params := synapse.ExtendedSQLPoolBlobAuditingPolicy{
		ExtendedSQLPoolBlobAuditingPolicyProperties: &synapse.ExtendedSQLPoolBlobAuditingPolicyProperties{
			State: synapse.BlobAuditingPolicyStateDisabled,
		},
	}

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, params); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
