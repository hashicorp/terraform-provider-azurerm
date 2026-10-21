// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func eventhubSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateEventHubName(),
		},

		"namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"partition_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validate.ValidateEventHubPartitionCount,
		},

		"message_retention": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.ValidateEventHubMessageRetentionCount,
			ExactlyOneOf: []string{"retention_description", "message_retention"},
		},

		"retention_description": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"cleanup_policy": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(eventhubs.PossibleValuesForCleanupPolicyRetentionDescription(), false),
					},

					"retention_time_in_hours": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ExactlyOneOf: []string{"retention_description.0.tombstone_retention_time_in_hours", "retention_description.0.retention_time_in_hours"},
					},

					"tombstone_retention_time_in_hours": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ExactlyOneOf: []string{"retention_description.0.retention_time_in_hours", "retention_description.0.tombstone_retention_time_in_hours"},
					},
				},
			},
		},

		"capture_description": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"skip_empty_archives": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"encoding": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(eventhubs.PossibleValuesForEncodingCaptureDescription(), false),
					},
					"interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      300,
						ValidateFunc: validation.IntBetween(60, 900),
					},
					"size_limit_in_bytes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      314572800,
						ValidateFunc: validation.IntBetween(10485760, 524288000),
					},
					"destination": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"EventHubArchive.AzureBlockBlob",
										// TODO: support `EventHubArchive.AzureDataLake` once supported in the Swagger / SDK
										// https://github.com/Azure/azure-rest-api-specs/issues/2255
										// BlobContainerName & StorageAccountID can then become Optional
									}, false),
								},
								"archive_name_format": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.ValidateEventHubArchiveNameFormat,
								},
								"blob_container_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"storage_account_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateStorageAccountID,
								},

								// Storage SAS is the default authentication type for capture destination, it's not supported by the API, so hard coding the value to align with azure portal behaviour.
								"storage_authentication_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "StorageSAS",
									ValidateFunc: validation.StringInSlice([]string{
										"StorageSAS",
										string(eventhubs.CaptureIdentityTypeSystemAssigned),
										string(eventhubs.CaptureIdentityTypeUserAssigned),
									}, false),
								},

								"storage_authentication_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: commonids.ValidateUserAssignedIdentityID,
								},
							},
						},
					},
				},
			},
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(eventhubs.EntityStatusActive),
			ValidateFunc: validation.StringInSlice([]string{
				string(eventhubs.EntityStatusActive),
				string(eventhubs.EntityStatusDisabled),
				string(eventhubs.EntityStatusSendDisabled),
			}, false),
		},

		"partition_ids": {
			Type:     pluginsdk.TypeSet,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:      pluginsdk.HashString,
			Computed: true,
		},
	}
}
