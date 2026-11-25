// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_custom_domain

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_app_environment_id": {
			Description:  "The ID of the Container Apps Managed Environment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
		},

		"certificate_blob_base64": {
			Description:  "The bundle of Private Key and Certificate for the Custom DNS Suffix as a base64 encoded PFX or PEM.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"certificate_password": {
			Description: "The password for the Certificate bundle.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			Sensitive:   true,
		},

		"dns_suffix": {
			Description:  "Custom DNS Suffix for the Container App Environment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
