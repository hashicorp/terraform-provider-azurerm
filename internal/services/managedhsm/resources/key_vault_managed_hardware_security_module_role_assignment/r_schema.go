// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_role_assignment

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_hsm_id": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			ValidateFunc: managedhsms.ValidateManagedHSMID,
			Required:     true,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"scope": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(/|/keys|/keys/.+)$`), "scope should be one of `/`, `/keys', `/keys/<key_name>`"),
		},

		"role_definition_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: roledefinitions.ValidateScopedRoleDefinitionID,
		},

		"principal_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
