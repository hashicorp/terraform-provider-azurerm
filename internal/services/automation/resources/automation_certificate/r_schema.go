// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Certificate.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Description:  "The name of the Automation Account in which the Certificate is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Description: "The description of this Automation Certificate.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"base64": {
			Description:  "Base64 encoded value of the certificate.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsBase64,
		},

		"exportable": {
			Description: "The is exportable flag of the certificate.",
			Type:        pluginsdk.TypeBool,
			Default:     false,
			Optional:    true,
		},

		"thumbprint": {
			Description: "The thumbprint for the certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
