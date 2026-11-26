// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_security_alert_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedserversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlManagedInstanceSecurityAlertPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstancesSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	disabledPolicy := managedserversecurityalertpolicies.ManagedServerSecurityAlertPolicy{
		Properties: &managedserversecurityalertpolicies.SecurityAlertsPolicyProperties{
			State: managedserversecurityalertpolicies.SecurityAlertsPolicyStateDisabled,
		},
	}

	if err = client.CreateOrUpdateThenPoll(ctx, managedInstanceId, disabledPolicy); err != nil {
		return fmt.Errorf("updating managed instance security alert policy: %v", err)
	}

	if _, err = client.Get(ctx, managedInstanceId); err != nil {
		return fmt.Errorf("deleting managed instance security alert policy: %v", err)
	}

	return nil
}
