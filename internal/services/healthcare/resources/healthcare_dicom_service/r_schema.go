// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_dicom_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func healthcareDicomServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Healthcare DICOM Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DicomServiceName(),
		},

		"workspace_id": {
			Description:  "The ID of the Healthcare Workspace where the Healthcare DICOM Service should exist.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"authentication": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authority": {
						Description: "The Azure Active Directory (tenant) that serves as the authentication authority to access the service. Authority must be registered to Azure AD and in the following format: <https://{Azure-AD-endpoint}/{tenant-id}.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"audience": {
						Description: "The intended audience to receive authentication tokens for the service.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
				},
			},
		},

		"private_endpoint": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"service_url": {
			Description: "The URL of the Healthcare DICOM Services.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"public_network_access_enabled": {
			Description: "Whether to enable public networks when data plane traffic coming from public networks while Private Endpoint is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"cors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C API sets defaults for these if omitted
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Description: "A list of allowed origins for CORS.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"allowed_headers": {
						Description: "A list of allowed headers for CORS.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"allowed_methods": {
						Description: "A list of allowed methods for CORS.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"max_age_in_seconds": {
						Description:  "The maximum age in seconds for the CORS configuration (must be between 0 and 99998 inclusive).",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 99998),
					},

					"allow_credentials": {
						Description: "Whether to allow credentials in CORS.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"data_partitions_enabled": {
			Description: "If data partitions are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"encryption_key_url": {
			Description:  "The URL of the key to use for encryption as part of the customer-managed key encryption settings. For more details, refer to the [Azure Customer-Managed Keys Overview](https://learn.microsoft.com/azure/storage/common/customer-managed-keys-overview).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"storage": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C API sets defaults for these if omitted
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"file_system_name": {
						Description:  "The filesystem name of connected Storage Account.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_account_id": {
						Description:  "The resource ID of connected Storage Account.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
