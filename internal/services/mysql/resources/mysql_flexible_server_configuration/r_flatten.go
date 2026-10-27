// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMySQLFlexibleServerConfigurationFlatten(d *pluginsdk.ResourceData, id *configurations.ConfigurationId, dbConfig *configurations.Configuration) error {
	d.Set("name", id.ConfigurationName)
	d.Set("server_name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	value := ""
	if dbConfig != nil {
		if props := dbConfig.Properties; props != nil {
			value = *props.Value
		}
	}
	d.Set("value", value)

	return pluginsdk.SetResourceIdentityData(d, id)
}
