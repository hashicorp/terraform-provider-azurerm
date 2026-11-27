// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_custom_hostname_binding

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_custom_hostname_binding/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name:hostname,site_name:app_service_name,resource_group_name" -test-name basicConfig -test-env-vars ARM_TEST_DNS_ZONE,ARM_TEST_DATA_RESOURCE_GROUP

const appServiceCustomHostnameBindingResourceName = "azurerm_app_service_custom_hostname_binding"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCustomHostnameBindingCreate,
		Read:   resourceAppServiceCustomHostnameBindingRead,
		Delete: resourceAppServiceCustomHostnameBindingDelete,
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&webapps.HostNameBindingId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&webapps.HostNameBindingId{}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AppServiceCustomHostnameBindingV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: appServiceCustomHostnameBindingSchema(),
	}
}
