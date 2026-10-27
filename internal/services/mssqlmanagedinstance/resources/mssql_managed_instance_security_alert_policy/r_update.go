// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_security_alert_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedserversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlManagedInstanceSecurityAlertPolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstancesSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	existing, err := client.Get(ctx, managedInstanceId)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	payload := existing.Model

	if d.HasChange("disabled_alerts") {
		if v, ok := d.GetOk("disabled_alerts"); ok {
			disabledAlerts := make([]string, 0)
			for _, v := range v.(*pluginsdk.Set).List() {
				disabledAlerts = append(disabledAlerts, v.(string))
			}
			payload.Properties.DisabledAlerts = &disabledAlerts
		} else {
			payload.Properties.DisabledAlerts = nil
		}
	}
	if d.HasChange("email_account_admins_enabled") {
		payload.Properties.EmailAccountAdmins = pointer.To(d.Get("email_account_admins_enabled").(bool))
	}

	if d.HasChange("retention_days") {
		payload.Properties.RetentionDays = pointer.To(int64(d.Get("retention_days").(int)))
	}

	if d.HasChange("email_addresses") {
		if v, ok := d.GetOk("email_addresses"); ok {
			emailAddresses := make([]string, 0)
			for _, v := range v.(*pluginsdk.Set).List() {
				emailAddresses = append(emailAddresses, v.(string))
			}
			payload.Properties.EmailAddresses = &emailAddresses
		} else {
			payload.Properties.EmailAddresses = nil
		}
	}

	if d.HasChange("enabled") {
		if d.Get("enabled").(bool) {
			payload.Properties.State = managedserversecurityalertpolicies.SecurityAlertsPolicyStateEnabled
		} else {
			payload.Properties.State = managedserversecurityalertpolicies.SecurityAlertsPolicyStateDisabled
		}
	}

	if d.HasChange("storage_account_access_key") {
		payload.Properties.StorageAccountAccessKey = pointer.To(d.Get("storage_account_access_key").(string))
	}

	// StorageAccountAccessKey cannot be passed in if it is empty. The api returns this as empty so we need to nil it before sending it back to the api
	if payload.Properties.StorageAccountAccessKey != nil && *payload.Properties.StorageAccountAccessKey == "" {
		payload.Properties.StorageAccountAccessKey = nil
	}

	if d.HasChange("storage_endpoint") {
		payload.Properties.StorageEndpoint = pointer.To(d.Get("storage_endpoint").(string))
	}

	// StorageEndpoint cannot be passed in if it is empty. The api returns this as empty so we need to nil it before sending it back to the api
	if payload.Properties.StorageEndpoint != nil && *payload.Properties.StorageEndpoint == "" {
		payload.Properties.StorageEndpoint = nil
	}

	if err = client.CreateOrUpdateThenPoll(ctx, managedInstanceId, *payload); err != nil {
		return fmt.Errorf("updating managed instance security alert policy: %v", err)
	}

	d.SetId(id.ID())

	return resourceMsSqlManagedInstanceSecurityAlertPolicyRead(d, meta)
}
