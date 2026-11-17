// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_active_directory_administrator

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/administratormicrosoftentras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var postgresqlFlexibleServerResourceName = "azurerm_postgresql_flexible_server"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgresqlFlexibleServerAdministratorCreate,
		Read:   resourcePostgresqlFlexibleServerAdministratorRead,
		Delete: resourcePostgresqlFlexibleServerAdministratorDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := administratormicrosoftentras.ParseAdministratorID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: postgresqlFlexibleServerActiveDirectoryAdministratorSchema(),
	}
}
