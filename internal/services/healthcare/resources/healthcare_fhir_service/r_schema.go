// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_fhir_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func healthcareFhirServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FhirServiceName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"location": commonschema.Location(),

		"kind": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(fhirservices.FhirServiceKindFhirNegativeRFour),
			ValidateFunc: validation.StringInSlice(fhirservices.PossibleValuesForFhirServiceKind(), false),
		},

		"access_policy_object_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"authentication": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authority": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"audience": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"smart_proxy_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		// can't use the registry ID due to the ID cannot be obtained when setting the property in state file
		"container_registry_login_server_url": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"oci_artifact": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"login_server": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"image_name": {
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
						Optional:     true,
					},

					"digest": {
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
						Optional:     true,
					},
				},
			},
		},

		"cors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"allowed_headers": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"allowed_methods": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"DELETE",
								"GET",
								"HEAD",
								"MERGE",
								"POST",
								"OPTIONS",
								"PUT",
								"PATCH",
							}, false),
						},
					},

					"max_age_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 2000000000),
					},

					"credentials_allowed": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"configuration_export_storage_account_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
