// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package user_assigned_identity

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/identities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),

		"isolation_scope": {
			Optional: true,
			Type:     pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				// `None` is not exposed
				string(identities.IsolationScopeRegional),
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"client_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"principal_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"tenant_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
