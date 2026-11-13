// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_configuration

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/configurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "flexible_server_name:server_id,resource_group_name:server_id" -test-params "backslash_quote,on"

var postgresqlFlexibleServerResourceName = "azurerm_postgresql_flexible_server"

const azurePostgresqlFlexibleServerConfigurationResourceName = "azurerm_postgresql_flexible_server_configuration"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFlexibleServerConfigurationCreateUpdate,
		Read:   resourceFlexibleServerConfigurationRead,
		Update: resourceFlexibleServerConfigurationCreateUpdate,
		Delete: resourceFlexibleServerConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&configurations.ConfigurationId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&configurations.ConfigurationId{}),

		Schema: postgresqlFlexibleServerConfigurationSchema(),
	}
}
