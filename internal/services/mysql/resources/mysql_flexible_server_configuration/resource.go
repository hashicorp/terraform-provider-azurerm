// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_configuration

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/configurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name mysql_flexible_server_configuration -service-package-name mysql -test-name characterSetServer -properties "name,resource_group_name,flexible_server_name:server_name" -known-values "subscription_id:data.Subscriptions.Primary"

var mysqlFlexibleServerResourceName = "azurerm_mysql_flexible_server"

var mysqlFlexibleServerConfigurationResourceName = "azurerm_mysql_flexible_server_configuration"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySQLFlexibleServerConfigurationCreate,
		Read:   resourceMySQLFlexibleServerConfigurationRead,
		Update: resourceMySQLFlexibleServerConfigurationUpdate,
		Delete: resourceMySQLFlexibleServerConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&configurations.ConfigurationId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&configurations.ConfigurationId{}),
		},

		Schema: mysqlFlexibleServerConfigurationSchema(),
	}
}
