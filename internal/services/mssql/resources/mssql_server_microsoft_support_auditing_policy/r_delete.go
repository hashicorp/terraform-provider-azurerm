// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_microsoft_support_auditing_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverdevopsaudit"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerMicrosoftSupportAuditingPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServerDevOpsAuditSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerMicrosoftSupportAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	params := serverdevopsaudit.ServerDevOpsAuditingSettings{
		Properties: &serverdevopsaudit.ServerDevOpsAuditSettingsProperties{
			State: serverdevopsaudit.BlobAuditingPolicyStateDisabled,
		},
	}

	if err = client.SettingsCreateOrUpdateThenPoll(ctx, serverId, params); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
