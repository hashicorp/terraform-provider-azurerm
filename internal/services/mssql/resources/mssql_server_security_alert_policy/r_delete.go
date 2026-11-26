// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_security_alert_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerSecurityAlertPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	disabledPolicy := serversecurityalertpolicies.ServerSecurityAlertPolicy{
		Properties: &serversecurityalertpolicies.SecurityAlertsPolicyProperties{
			State: serversecurityalertpolicies.SecurityAlertsPolicyStateDisabled,
		},
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	if err = client.CreateOrUpdateThenPoll(ctx, serverId, disabledPolicy); err != nil {
		return fmt.Errorf("updating mssql server security alert policy: %+v", err)
	}

	if _, err = client.Get(ctx, serverId); err != nil {
		return fmt.Errorf("deleting mssql server security alert policy: %+v", err)
	}

	return nil
}
