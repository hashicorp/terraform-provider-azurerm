// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-expect-non-empty

var mysqlFlexibleServerResourceName = "azurerm_mysql_flexible_server"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMysqlFlexibleServerCreate,
		Read:   resourceMysqlFlexibleServerRead,
		Update: resourceMysqlFlexibleServerUpdate,
		Delete: resourceMysqlFlexibleServerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(2 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(2 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&servers.FlexibleServerId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&servers.FlexibleServerId{}),
		},

		Schema: mysqlFlexibleServerSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("storage.0.size_gb", func(ctx context.Context, old, new, meta any) bool {
				return new.(int) < old.(int)
			}),
		),
	}
}
