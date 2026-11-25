// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_configuration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMySQLFlexibleServerConfigurationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Configurations
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// NOTE: this resource intentionally doesn't support Requires Import
	//       since a fallback route is created by default

	id := configurations.NewConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))

	locks.ByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)
	defer locks.UnlockByName(id.FlexibleServerName, mysqlFlexibleServerResourceName)

	updateCallback := sdk.SetIDAndIdentityCallback(meta, &id, d)

	if strings.EqualFold(id.ConfigurationName, "gtid_mode") {
		if err := mysqlFlexibleServerConfigurationUpdateGITDMode(ctx, client, id, d.Get("value").(string), updateCallback); err != nil {
			return fmt.Errorf("creating %s: %v", id, err)
		}
	} else {
		payload := configurations.Configuration{
			Properties: &configurations.ConfigurationProperties{
				Value: pointer.To(d.Get("value").(string)),
			},
		}

		if err := client.UpdateCallbackThenPoll(ctx, id, payload, updateCallback); err != nil {
			return fmt.Errorf("creating %s: %v", id, err)
		}
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceMySQLFlexibleServerConfigurationRead(d, meta)
}
