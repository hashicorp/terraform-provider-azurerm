// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_security_alert_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerSecurityAlertPolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	policy, err := client.Get(ctx, serverId)
	if err != nil {
		return fmt.Errorf("retrieving mssql server security alert policy %s: %+v", id, err)
	}

	model := policy.Model
	if model == nil {
		return fmt.Errorf("retrieving mssql server security alert policy %s: model was nil", id)
	}

	if model.Name == nil {
		return fmt.Errorf("reading mssql server security alert policy %s: name was nil", id)
	}

	if model.Properties == nil {
		return fmt.Errorf("reading mssql server security alert policy %s: properties was nil", id)
	}

	payload := serversecurityalertpolicies.ServerSecurityAlertPolicy{}
	props := model.Properties

	if d.HasChange("state") {
		props.State = serversecurityalertpolicies.SecurityAlertsPolicyState(d.Get("state").(string))
	}

	if d.HasChange("disabled_alerts") {
		disabledAlerts := make([]string, 0)
		if v, ok := d.GetOk("disabled_alerts"); ok {
			for _, v := range v.(*pluginsdk.Set).List() {
				disabledAlerts = append(disabledAlerts, v.(string))
			}
		}
		props.DisabledAlerts = pointer.To(disabledAlerts)
	}

	if d.HasChange("email_addresses") {
		emailAddresses := make([]string, 0)
		if v, ok := d.GetOk("email_addresses"); ok {
			for _, v := range v.(*pluginsdk.Set).List() {
				emailAddresses = append(emailAddresses, v.(string))
			}
		}
		props.EmailAddresses = pointer.To(emailAddresses)
	}

	if d.HasChange("email_account_admins_enabled") {
		props.EmailAccountAdmins = pointer.To(d.Get("email_account_admins_enabled").(bool))
	}

	if d.HasChange("retention_days") {
		var retentionDays *int64
		if v, ok := d.GetOk("retention_days"); ok {
			retentionDays = pointer.To(int64(v.(int)))
		}
		props.RetentionDays = retentionDays
	}

	if d.HasChange("storage_endpoint") {
		props.StorageEndpoint = nil
		if v, ok := d.GetOk("storage_endpoint"); ok {
			props.StorageEndpoint = pointer.To(v.(string))
		}
	}

	// NOTE: 'storage_account_access_key' it is not returned
	// by the API, so we need to get it from state...
	props.StorageAccountAccessKey = nil
	if v, ok := d.GetOk("storage_account_access_key"); ok {
		props.StorageAccountAccessKey = pointer.To(v.(string))
	}

	payload.Properties = props

	if err = client.CreateOrUpdateThenPoll(ctx, serverId, payload); err != nil {
		return fmt.Errorf("updating mssql server security alert policy: %+v", err)
	}

	return resourceMsSqlServerSecurityAlertPolicyRead(d, meta)
}
