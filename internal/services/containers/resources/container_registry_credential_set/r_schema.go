// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_credential_set

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Container Registry Credential Set.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"container_registry_id": {
			Description:  "The ID of the Container Registry.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},
		"login_server": {
			Description: "The login server for the Credential Set.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"authentication_credentials": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"username_secret_id": {
						Description:  "The URI of the secret containing the username in a Key Vault.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeSecret),
					},
					"password_secret_id": {
						Description:  "The URI of the secret containing the password in a Key Vault.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeSecret),
					},
				},
			},
		},
		// This property relies on a pandora workaround due to a swagger issue
		// https://github.com/Azure/azure-rest-api-specs/issues/32154
		"identity": commonschema.SystemAssignedIdentityRequired(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
