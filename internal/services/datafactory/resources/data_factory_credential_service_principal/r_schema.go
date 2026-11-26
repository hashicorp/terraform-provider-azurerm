// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_credential_service_principal

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Credential.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_factory_id": {
			Description:  "The Data Factory ID in which to associate the Credential with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: credentials.ValidateFactoryID,
		},

		"tenant_id": {
			Description:  "The Tenant ID of the Service Principal.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"service_principal_id": {
			Description:  "The Client ID of the Service Principal.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"service_principal_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Description:  "The name of the Linked Service to use for the Service Principal Key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Description:  "The name of the Secret in the Key Vault.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_version": {
						Description:  "The version of the Secret in the Key Vault.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"description": {
			Description:  "The description for the Data Factory Credential.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"annotations": { // this property is not visible in the azure portal
			Description: "List of tags that can be used for describing the Data Factory Credential.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
