// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package portal_tenant_configuration

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePortalTenantConfigurationCreateUpdate,
		Read:   resourcePortalTenantConfigurationRead,
		Update: resourcePortalTenantConfigurationCreateUpdate,
		Delete: resourcePortalTenantConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: schema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PortalTenantConfigurationID(id)
			return err
		}),

		Schema: portalTenantConfigurationSchema(),
	}
}
