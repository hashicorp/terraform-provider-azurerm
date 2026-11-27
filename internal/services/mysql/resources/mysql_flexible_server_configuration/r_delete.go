// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	mysqlClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMySQLFlexibleServerConfigurationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := mysqlClient.From(meta).FlexibleServers.Configurations
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)

	// "delete" = resetting this to the default value
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := configurations.Configuration{
		Properties: &configurations.ConfigurationProperties{
			// we can alternatively set `source: "system-default"`
			Value: resp.Model.Properties.DefaultValue,
		},
	}

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("resetting %s to it's default value: %+v", *id, err)
	}

	return nil
}
