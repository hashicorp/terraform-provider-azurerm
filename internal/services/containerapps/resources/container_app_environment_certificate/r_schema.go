// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Container Apps Environment Certificate.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateName,
		},

		"container_app_environment_id": {
			Description:  "The Container App Managed Environment ID to configure this Certificate on.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificates.ValidateManagedEnvironmentID,
		},

		"certificate_blob_base64": {
			Description:  "The Certificate Private Key as a base64 encoded PFX or PEM.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsBase64,
			RequiredWith: []string{"certificate_password"},
			ExactlyOneOf: []string{"certificate_key_vault", "certificate_blob_base64"},
		},

		"certificate_password": {
			Description:  "The password for the Certificate.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			RequiredWith: []string{"certificate_blob_base64"},
		},

		"certificate_key_vault": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"certificate_key_vault", "certificate_blob_base64"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity": {
						Description: "The Managed Identity to authenticate with Azure Key Vault.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
						Default:     "System",
						ValidateFunc: validation.Any(
							validation.StringInSlice([]string{"System"}, false),
							commonids.ValidateUserAssignedIdentityID,
						),
					},
					"key_vault_secret_id": {
						Description:  "The ID of the Key Vault Secret containing the certificate.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeSecret),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subject_name": {
			Description: "The Subject Name for the Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"issuer": {
			Description: "The Certificate Issuer.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"issue_date": {
			Description: "The date of issue for the Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"expiration_date": {
			Description: "The expiration date for the Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"thumbprint": {
			Description: "The Thumbprint of the Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
