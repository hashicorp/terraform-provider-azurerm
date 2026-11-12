// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	namespaceName := ""
	resourceGroupName := ""
	if v := d.Get("namespace_id").(string); v != "" {
		namespaceId, err := namespaces.ParseNamespaceID(v)
		if err != nil {
			return err
		}
		namespaceName = namespaceId.NamespaceName
		resourceGroupName = namespaceId.ResourceGroupName
	}

	id := eventhubs.NewEventhubID(subscriptionId, resourceGroupName, namespaceName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventhub", id.ID())
		}
	}

	eventhubStatus := eventhubs.EntityStatus(d.Get("status").(string))
	parameters := eventhubs.Eventhub{
		Properties: &eventhubs.EventhubProperties{
			PartitionCount: pointer.To(int64(d.Get("partition_count").(int))),
			Status:         &eventhubStatus,
		},
	}

	if _, ok := d.GetOk("retention_description"); ok {
		parameters.Properties.RetentionDescription = expandEventHubRetentionDescription(d)
	}

	if _, ok := d.GetOk("message_retention"); ok {
		parameters.Properties.MessageRetentionInDays = pointer.To(int64(d.Get("message_retention").(int)))
	}

	if _, ok := d.GetOk("capture_description"); ok {
		parameters.Properties.CaptureDescription = expandEventHubCaptureDescription(d)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceEventHubRead(d, meta)
}

func expandEventHubRetentionDescription(d *pluginsdk.ResourceData) *eventhubs.RetentionDescription {
	inputs := d.Get("retention_description").([]any)
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}
	input := inputs[0].(map[string]any)

	cleanupPolicy := input["cleanup_policy"].(string)

	retentionDescription := eventhubs.RetentionDescription{
		CleanupPolicy: pointer.ToEnum[eventhubs.CleanupPolicyRetentionDescription](cleanupPolicy),
	}

	if cleanupPolicy == string(eventhubs.CleanupPolicyRetentionDescriptionDelete) {
		retentionTimeInHours := input["retention_time_in_hours"].(int)
		retentionDescription.RetentionTimeInHours = pointer.To(int64(retentionTimeInHours))
	} else {
		tombstoneRetentionTimeInHours := input["tombstone_retention_time_in_hours"].(int)
		retentionDescription.TombstoneRetentionTimeInHours = pointer.To(int64(tombstoneRetentionTimeInHours))
	}

	return pointer.To(retentionDescription)
}

func expandEventHubCaptureDescription(d *pluginsdk.ResourceData) *eventhubs.CaptureDescription {
	inputs := d.Get("capture_description").([]any)
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}
	input := inputs[0].(map[string]any)

	enabled := input["enabled"].(bool)
	encoding := input["encoding"].(string)
	intervalInSeconds := input["interval_in_seconds"].(int)
	sizeLimitInBytes := input["size_limit_in_bytes"].(int)
	skipEmptyArchives := input["skip_empty_archives"].(bool)

	captureDescription := eventhubs.CaptureDescription{
		Enabled: pointer.To(enabled),
		Encoding: func() *eventhubs.EncodingCaptureDescription {
			return pointer.ToEnum[eventhubs.EncodingCaptureDescription](encoding)
		}(),
		IntervalInSeconds: pointer.To(int64(intervalInSeconds)),
		SizeLimitInBytes:  pointer.To(int64(sizeLimitInBytes)),
		SkipEmptyArchives: pointer.To(skipEmptyArchives),
	}

	if v, ok := input["destination"]; ok {
		destinations := v.([]any)
		if len(destinations) > 0 {
			destination := destinations[0].(map[string]any)

			destinationName := destination["name"].(string)
			archiveNameFormat := destination["archive_name_format"].(string)
			blobContainerName := destination["blob_container_name"].(string)
			storageAccountId := destination["storage_account_id"].(string)

			captureDescription.Destination = &eventhubs.Destination{
				Name: pointer.To(destinationName),
				Properties: &eventhubs.DestinationProperties{
					ArchiveNameFormat:        pointer.To(archiveNameFormat),
					BlobContainer:            pointer.To(blobContainerName),
					StorageAccountResourceId: pointer.To(storageAccountId),
				},
			}

			if destinationAuthType := destination["storage_authentication_type"]; destinationAuthType != nil && destinationAuthType.(string) != "" {
				authType := eventhubs.CaptureIdentityType(destinationAuthType.(string))
				if authType != "StorageSAS" {
					captureDescription.Destination.Identity = &eventhubs.CaptureIdentity{
						Type: pointer.To(authType),
					}
				}
			}

			if destinationAuthTypeId := destination["storage_authentication_id"]; destinationAuthTypeId != nil && destinationAuthTypeId.(string) != "" {
				captureDescription.Destination.Identity.UserAssignedIdentity = pointer.To(destinationAuthTypeId.(string))
			}
		}
	}

	return &captureDescription
}
