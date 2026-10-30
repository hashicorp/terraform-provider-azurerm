// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourcePostgresqlFlexibleServerConfigurationFlatten(d *pluginsdk.ResourceData, id *configurations.ConfigurationId, model *configurations.Configuration) error {
	d.Set("name", id.ConfigurationName)
	d.Set("server_id", configurations.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName).ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("value", props.Value)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
