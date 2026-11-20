// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_security_alert_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedserversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlManagedInstanceSecurityAlertPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstancesSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	result, err := client.Get(ctx, managedInstanceId)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			log.Printf("[WARN] managed instance security alert policy %v not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request to managed instance security alert policy: %+v", err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("managed_instance_name", id.ManagedInstanceName)

	if result.Model != nil {
		if props := result.Model.Properties; props != nil {
			d.Set("enabled", props.State == managedserversecurityalertpolicies.SecurityAlertsPolicyStateEnabled)

			if props.DisabledAlerts != nil {
				disabledAlerts := pluginsdk.NewSet(pluginsdk.HashString, []any{})
				for _, v := range *props.DisabledAlerts {
					if v != "" {
						disabledAlerts.Add(v)
					}
				}

				d.Set("disabled_alerts", disabledAlerts)
			}

			if props.EmailAccountAdmins != nil {
				d.Set("email_account_admins_enabled", props.EmailAccountAdmins)
			}

			if props.EmailAddresses != nil {
				emailAddresses := pluginsdk.NewSet(pluginsdk.HashString, []any{})
				for _, v := range *props.EmailAddresses {
					if v != "" {
						emailAddresses.Add(v)
					}
				}

				d.Set("email_addresses", emailAddresses)
			}

			if props.RetentionDays != nil {
				d.Set("retention_days", int(*props.RetentionDays))
			}

			if v, ok := d.GetOk("storage_account_access_key"); ok {
				d.Set("storage_account_access_key", v)
			}

			if props.StorageEndpoint != nil {
				d.Set("storage_endpoint", props.StorageEndpoint)
			}
		}
	}

	return nil
}
