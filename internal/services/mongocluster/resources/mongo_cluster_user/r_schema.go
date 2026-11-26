// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster_user

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/users"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// `object_id` is actually resource name
		"object_id": {
			Description:  "The Object ID of the Entra ID User or Service Principal.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"mongo_cluster_id": commonschema.ResourceIDReferenceRequiredForceNew(&users.MongoClusterId{}, "Mongo Cluster"),

		"identity_provider_type": {
			Description: "The identity provider type for the Mongo Cluster User.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice(
				users.PossibleValuesForIdentityProviderType(),
				false,
			),
		},

		"principal_type": {
			Description: "The principal type for the Mongo Cluster User.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice(
				users.PossibleValuesForEntraPrincipalType(),
				false,
			),
		},

		"role": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"database": {
						Description:  "The name of the database to which the role will be applied.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"name": {
						Description: "The role name.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
						ValidateFunc: validation.StringInSlice(
							users.PossibleValuesForUserRole(),
							false,
						),
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
