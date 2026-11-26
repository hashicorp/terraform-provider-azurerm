// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_security_alert_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerSecurityAlertPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerSecurityAlertPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	result, err := client.Get(ctx, serverId)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			log.Printf("[WARN] mssql server security alert policy %s: not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request to mssql server security alert policy: %+v", err)
	}

	model := result.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("retrieving %s: Properties was nil", id)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)
	d.Set("state", string(props.State))

	disabledAlerts := pluginsdk.NewSet(pluginsdk.HashString, []any{})
	if props.DisabledAlerts != nil {
		for _, v := range *props.DisabledAlerts {
			if v != "" {
				disabledAlerts.Add(v)
			}
		}
	}
	d.Set("disabled_alerts", disabledAlerts)

	d.Set("email_account_admins_enabled", props.EmailAccountAdmins)

	emailAddresses := pluginsdk.NewSet(pluginsdk.HashString, []any{})
	if props.EmailAddresses != nil {
		for _, v := range *props.EmailAddresses {
			if v != "" {
				emailAddresses.Add(v)
			}
		}
	}
	d.Set("email_addresses", emailAddresses)

	var retentionDays int
	if props.RetentionDays != nil {
		retentionDays = int(*props.RetentionDays)
	}
	d.Set("retention_days", retentionDays)

	d.Set("storage_endpoint", pointer.From(props.StorageEndpoint))

	// NOTE: 'storage_account_access_key' field is not returned by the API
	// so we need to pull it from the state...
	var accessKey string
	if v, ok := d.GetOk("storage_account_access_key"); ok {
		accessKey = v.(string)
	}
	d.Set("storage_account_access_key", accessKey)

	return nil
}
