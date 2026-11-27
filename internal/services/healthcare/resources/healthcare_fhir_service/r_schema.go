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
			Description:  "The name of the Healthcare FHIR Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FhirServiceName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_id": {
			Description:  "The ID of the Healthcare Workspace where the Healthcare FHIR Service should exist.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"location": commonschema.Location(),

		"kind": {
			Description:  "The kind of the Healthcare FHIR Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(fhirservices.FhirServiceKindFhirNegativeRFour),
			ValidateFunc: validation.StringInSlice(fhirservices.PossibleValuesForFhirServiceKind(), false),
		},

		"access_policy_object_ids": {
			Description: "A list of the access policies of the service instance.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
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
						Description:  "The Azure Active Directory (tenant) that serves as the authentication authority to access the service. Authority must be registered to Azure AD and in the following format: <https://{Azure-AD-endpoint}/{tenant-id}.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"audience": {
						Description: "The intended audience to receive authentication tokens for the service.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},

					"smart_proxy_enabled": {
						Description: "Whether smart proxy is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		// can't use the registry ID due to the ID cannot be obtained when setting the property in state file
		"container_registry_login_server_url": {
			Description: "A list of azure container registry settings used for convert data operation of the service instance.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
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
						Description:  "An Azure container registry used for export operations of the service instance.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"image_name": {
						Description:  "An image within Azure container registry used for export operations of the service instance.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"digest": {
						Description:  "A digest of an image within Azure container registry used for export operations of the service instance to narrow the artifacts down.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
						Description: "A set of origins to be allowed via CORS.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						MaxItems:    64,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"allowed_headers": {
						Description: "A set of headers to be allowed via CORS.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						MaxItems:    64,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"allowed_methods": {
						Description: "The methods to be allowed via CORS.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						MaxItems:    64,
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
						Description:  "The max age to be allowed via CORS.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 2000000000),
					},

					"credentials_allowed": {
						Description: "If credentials are allowed via CORS.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"configuration_export_storage_account_name": {
			Description:  "The name of the Storage Account which the operation configuration information is exported to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"public_network_access_enabled": {
			Description: "Whether public networks access is enabled.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
