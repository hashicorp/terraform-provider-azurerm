// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	mysqlClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMySQLFlexibleServerConfigurationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := mysqlClient.From(meta).FlexibleServers.Configurations
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	payload := configurations.Configuration{
		Properties: &configurations.ConfigurationProperties{
			Value: pointer.To(d.Get("value").(string)),
		},
	}

	// NOTE: this resource intentionally doesn't support Requires Import
	//       since a fallback route is created by default

	id := configurations.NewConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))

	locks.ByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)

	if err := client.UpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceMySQLFlexibleServerConfigurationRead(d, meta)
}
