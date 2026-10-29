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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateName,
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificates.ValidateManagedEnvironmentID,
		},

		"certificate_blob_base64": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsBase64,
			ExactlyOneOf: []string{"certificate_key_vault", "certificate_blob_base64"},
			RequiredWith: []string{"certificate_password"},
		},

		"certificate_password": {
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "System",
						ValidateFunc: validation.Any(
							validation.StringInSlice([]string{"System"}, false),
							commonids.ValidateUserAssignedIdentityID,
						),
					},
					"key_vault_secret_id": {
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
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Subject Name for the Certificate.",
		},

		"issuer": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Certificate Issuer.",
		},

		"issue_date": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The date of issue for the Certificate.",
		},

		"expiration_date": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The expiration date for the Certificate.",
		},

		"thumbprint": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Thumbprint of the Certificate.",
		},
	}
}
