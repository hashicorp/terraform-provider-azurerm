// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_database

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/databases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name mysql_flexible_database -service-package-name mysql -properties "name,resource_group_name,flexible_server_name:server_name" -known-values "subscription_id:data.Subscriptions.Primary"

var mysqlFlexibleDatabaseResourceName = "azurerm_mysql_flexible_database"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceMySqlFlexibleDatabaseCreate,
		Read:     resourceMySqlFlexibleDatabaseRead,
		Delete:   resourceMySqlFlexibleDatabaseDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&databases.DatabaseId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&databases.DatabaseId{}),
		},

		Schema: mysqlFlexibleDatabaseSchema(),
	}
}
