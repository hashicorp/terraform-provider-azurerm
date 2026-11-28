// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobcontainers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageContainerCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)
	accessLevelRaw := d.Get("container_access_type").(string)
	metaDataRaw := d.Get("metadata").(map[string]any)
	metaData := helpers.ExpandMetaData(metaDataRaw)

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := commonids.NewStorageContainerID(subscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, containerName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := containerClient.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(StorageContainerResourceName, id.ID())
		}
	}

	payload := blobcontainers.BlobContainer{
		Properties: &blobcontainers.ContainerProperties{
			PublicAccess: pointer.ToEnum[blobcontainers.PublicAccess](containerAccessTypeConversionMap[accessLevelRaw]),
			Metadata:     pointer.To(metaData),
		},
	}

	if encryptionScope := d.Get("default_encryption_scope"); encryptionScope.(string) != "" {
		payload.Properties.DefaultEncryptionScope = pointer.To(encryptionScope.(string))
		payload.Properties.DenyEncryptionScopeOverride = pointer.To(false)

		if encryptionScopeOverrideEnabled := d.Get("encryption_scope_override_enabled"); !encryptionScopeOverrideEnabled.(bool) {
			payload.Properties.DenyEncryptionScopeOverride = pointer.To(true)
		}
	}

	resp, err := containerClient.Create(ctx, id, payload)
	if err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	if meta.(*clients.Client).Features.PersistIDOnCreateBeforePollingForCompletion {
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	}

	pollerType := custompollers.NewStorageContainerCreatePoller(containerClient, id, resp.HttpResponse)
	poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	if err = poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for creation of %s: %v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceStorageContainerRead(d, meta)
}
