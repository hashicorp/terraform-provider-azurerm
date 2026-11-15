// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceEventHubFlatten(d *pluginsdk.ResourceData, id *eventhubs.EventhubId, model *eventhubs.Eventhub) error {
	d.Set("name", id.EventhubName)

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	d.Set("namespace_id", namespaceId.ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("partition_count", props.PartitionCount)
			d.Set("partition_ids", props.PartitionIds)
			d.Set("status", string(*props.Status))

			if err := d.Set("capture_description", flattenEventHubCaptureDescription(props.CaptureDescription)); err != nil {
				return err
			}

			// TODO - the `props.RetentionDescription.TombstoneRetentionTimeInHours == nil` check can be removed when https://github.com/Azure/azure-rest-api-specs/issues/36018 is fixed
			if props.RetentionDescription == nil || props.RetentionDescription.TombstoneRetentionTimeInHours == nil {
				d.Set("message_retention", props.MessageRetentionInDays)
			}
			if err := d.Set("retention_description", flattenEventHubRetentionDescription(props.RetentionDescription)); err != nil {
				return err
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenEventHubRetentionDescription(description *eventhubs.RetentionDescription) []any {
	results := make([]any, 0)

	if description != nil {
		output := make(map[string]any)

		if cleanupPolicy := description.CleanupPolicy; cleanupPolicy != nil {
			output["cleanup_policy"] = pointer.FromEnum(cleanupPolicy)
		}

		if retentionTimeInHours := description.RetentionTimeInHours; retentionTimeInHours != nil && output["cleanup_policy"] == "Delete" {
			output["retention_time_in_hours"] = pointer.From(retentionTimeInHours)
		}

		if tombstoneRetentionTimeInHours := description.TombstoneRetentionTimeInHours; tombstoneRetentionTimeInHours != nil && output["cleanup_policy"] == "Compact" {
			output["tombstone_retention_time_in_hours"] = pointer.From(tombstoneRetentionTimeInHours)
		}

		results = append(results, output)
	}

	return results
}

func flattenEventHubCaptureDescription(description *eventhubs.CaptureDescription) []any {
	results := make([]any, 0)

	if description != nil {
		output := make(map[string]any)

		if enabled := description.Enabled; enabled != nil {
			output["enabled"] = *enabled
		}

		if skipEmptyArchives := description.SkipEmptyArchives; skipEmptyArchives != nil {
			output["skip_empty_archives"] = *skipEmptyArchives
		}

		encoding := ""
		if description.Encoding != nil {
			encoding = string(*description.Encoding)
		}
		output["encoding"] = encoding

		if interval := description.IntervalInSeconds; interval != nil {
			output["interval_in_seconds"] = *interval
		}

		if size := description.SizeLimitInBytes; size != nil {
			output["size_limit_in_bytes"] = *size
		}

		if destination := description.Destination; destination != nil {
			destinationOutput := make(map[string]any)

			if name := destination.Name; name != nil {
				destinationOutput["name"] = *name
			}

			if props := destination.Properties; props != nil {
				if archiveNameFormat := props.ArchiveNameFormat; archiveNameFormat != nil {
					destinationOutput["archive_name_format"] = *archiveNameFormat
				}
				if blobContainerName := props.BlobContainer; blobContainerName != nil {
					destinationOutput["blob_container_name"] = *blobContainerName
				}
				if storageAccountId := props.StorageAccountResourceId; storageAccountId != nil {
					destinationOutput["storage_account_id"] = *storageAccountId
				}
			}

			destinationOutput["storage_authentication_type"] = "StorageSAS"
			if storageIdentity := destination.Identity; storageIdentity != nil {
				if storageAuthType := storageIdentity.Type; storageAuthType != nil {
					destinationOutput["storage_authentication_type"] = pointer.From(storageIdentity.Type)
				}
				destinationOutput["storage_authentication_id"] = pointer.From(storageIdentity.UserAssignedIdentity)
			}

			output["destination"] = []any{destinationOutput}
		}

		results = append(results, output)
	}

	return results
}
