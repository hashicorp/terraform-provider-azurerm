// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package portal_tenant_configuration

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func portalTenantConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"private_markdown_storage_enforced": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},
	}
}
