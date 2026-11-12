// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package tenant_template_deployment

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: tenantTemplateDeploymentResourceCreate,
		Read:   tenantTemplateDeploymentResourceRead,
		Update: tenantTemplateDeploymentResourceUpdate,
		Delete: tenantTemplateDeploymentResourceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.TenantTemplateDeploymentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(180 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(180 * time.Minute),
		},

		// (@jackofallops - lintignore needed as we need to make sure the JSON is usable in `output_content`)

		// lintignore:S033
		Schema: tenantTemplateDeploymentSchema(),
	}
}
