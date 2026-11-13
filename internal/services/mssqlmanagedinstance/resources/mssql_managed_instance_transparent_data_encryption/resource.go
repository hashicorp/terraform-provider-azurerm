// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_transparent_data_encryption

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlManagedInstanceTransparentDataEncryptionCreateUpdate,
		Read:   resourceMsSqlManagedInstanceTransparentDataEncryptionRead,
		Update: resourceMsSqlManagedInstanceTransparentDataEncryptionCreateUpdate,
		Delete: resourceMsSqlManagedInstanceTransparentDataEncryptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedInstanceEncryptionProtectorID(id)

			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: mssqlManagedInstanceTransparentDataEncryptionSchema(),
	}
}
