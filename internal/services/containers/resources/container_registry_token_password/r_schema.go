// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_token_password

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_registry_token_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: tokens.ValidateTokenID,
		},

		"password1": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiry": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// TODO: Remove the force new and add customize diff to SetNewComputed on the `value` once https://github.com/hashicorp/terraform-plugin-sdk/issues/459 is addressed.
						ForceNew:         true,
						ValidateFunc:     validation.IsRFC3339Time,
						DiffSuppressFunc: suppress.RFC3339Time,
					},

					"value": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},

		"password2": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiry": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// TODO: Remove the force new and add customize diff to SetNewComputed on the `value` once https://github.com/hashicorp/terraform-plugin-sdk/issues/459 is addressed.
						ForceNew:         true,
						ValidateFunc:     validation.IsRFC3339Time,
						DiffSuppressFunc: suppress.RFC3339Time,
					},

					"value": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
