// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func automationCertificateSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"base64": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsBase64,
		},

		"exportable": {
			Type:     pluginsdk.TypeBool,
			Default:  false,
			Optional: true,
		},

		"thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
