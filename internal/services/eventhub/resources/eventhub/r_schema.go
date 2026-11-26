// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the EventHub resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateEventHubName(),
		},

		"namespace_id": {
			Description:  "The ID of the EventHub Namespace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"partition_count": {
			Description:  "The current number of shards on the Event Hub.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validate.ValidateEventHubPartitionCount,
		},

		"message_retention": {
			Description:  "The number of days to retain the events for this Event Hub.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.ValidateEventHubMessageRetentionCount,
			ExactlyOneOf: []string{"retention_description", "message_retention"},
		},

		"retention_description": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"cleanup_policy": {
						Description:  "The Cleanup Policy for the EventHub.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(eventhubs.PossibleValuesForCleanupPolicyRetentionDescription(), false),
					},

					"retention_time_in_hours": {
						Description:  "The number of hours to retain the events for this Event Hub. The value is only used when `cleanup_policy` is `Delete`.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ExactlyOneOf: []string{"retention_description.0.tombstone_retention_time_in_hours", "retention_description.0.retention_time_in_hours"},
					},

					"tombstone_retention_time_in_hours": {
						Description:  "The number of hours to retain the tombstones markers of a compacted Event Hub. The value is only used when `cleanup_policy` is `Compact`.",
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
						Description: "Whether the Capture Description is enabled.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"skip_empty_archives": {
						Description: "Whether empty files should not be emitted if no events occur during the Capture time window.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"encoding": {
						Description:  "The Encoding used for the Capture Description.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(eventhubs.PossibleValuesForEncodingCaptureDescription(), false),
					},
					"interval_in_seconds": {
						Description:  "The time interval in seconds at which the capture will happen.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      300,
						ValidateFunc: validation.IntBetween(60, 900),
					},
					"size_limit_in_bytes": {
						Description:  "The amount of data built up in your EventHub before a Capture Operation occurs.",
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
									Description: "The Name of the Destination where the capture should take place.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"EventHubArchive.AzureBlockBlob",
										// TODO: support `EventHubArchive.AzureDataLake` once supported in the Swagger / SDK
										// https://github.com/Azure/azure-rest-api-specs/issues/2255
										// BlobContainerName & StorageAccountID can then become Optional
									}, false),
								},
								"archive_name_format": {
									Description:  "The Blob naming convention for archiving. e.g. `{Namespace}/{EventHub}/{PartitionId}/{Year}/{Month}/{Day}/{Hour}/{Minute}/{Second}`. Here all the parameters (Namespace, EventHub, etc.) are mandatory irrespective of order.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.ValidateEventHubArchiveNameFormat,
								},
								"blob_container_name": {
									Description: "The name of the Container within the Blob Storage Account where messages should be archived.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"storage_account_id": {
									Description:  "The ID of the Blob Storage Account where messages should be archived.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateStorageAccountID,
								},

								// Storage SAS is the default authentication type for capture destination, it's not supported by the API, so hard coding the value to align with azure portal behaviour.
								"storage_authentication_type": {
									Description: "The identity used to authenticate the Blob Storage Account where messages should be archived.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									Default:     "StorageSAS",
									ValidateFunc: validation.StringInSlice([]string{
										"StorageSAS",
										string(eventhubs.CaptureIdentityTypeSystemAssigned),
										string(eventhubs.CaptureIdentityTypeUserAssigned),
									}, false),
								},

								"storage_authentication_id": {
									Description:  "The ID of the User Assigned Identity that is used to authenticate the Blob Storage Account where messages should be archived.",
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
			Description: "The status of the Event Hub resource.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(eventhubs.EntityStatusActive),
			ValidateFunc: validation.StringInSlice([]string{
				string(eventhubs.EntityStatusActive),
				string(eventhubs.EntityStatusDisabled),
				string(eventhubs.EntityStatusSendDisabled),
			}, false),
		},

		"partition_ids": {
			Description: "The identifiers for partitions created for Event Hubs.",
			Type:        pluginsdk.TypeSet,
			Computed:    true,
			Set:         pluginsdk.HashString,
			Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
	}
}
