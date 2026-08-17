// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobcontainers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

var containerAccessTypeConversionMap = map[string]string{
	"blob":      "Blob",
	"container": "Container",
	"private":   "None",
	"Blob":      "blob",
	"Container": "container",
	"None":      "private",
	"":          "private",
}

func resourceStorageContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageContainerCreate,
		Read:   resourceStorageContainerRead,
		Delete: resourceStorageContainerDelete,
		Update: resourceStorageContainerUpdate,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := commonids.ParseStorageContainerID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ContainerV0ToV1{},
			1: migration.StorageContainerV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageContainerName,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"container_access_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "private",
				ValidateFunc: validation.StringInSlice([]string{
					string(containers.Blob),
					string(containers.Container),
					"private",
				}, false),
			},

			"default_encryption_scope": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true, // needed because a dummy value is returned when unspecified
				ForceNew:     true,
				ValidateFunc: validate.StorageEncryptionScopeName,
			},

			"encryption_scope_override_enabled": {
				Type:         pluginsdk.TypeBool,
				Optional:     true,
				Default:      true, // defaulting to false would be preferable here, but the API defaults this to true when unspecified
				ForceNew:     true,
				RequiredWith: []string{"default_encryption_scope"},
			},

			"metadata": MetaDataComputedSchema(),

			"has_immutability_policy": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"has_legal_hold": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageContainerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)
	accessLevelRaw := d.Get("container_access_type").(string)
	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

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
			return tf.ImportAsExistsError("azurerm_storage_container", id.ID())
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

	pollerType := custompollers.NewStorageContainerCreatePoller(containerClient, id, resp.HttpResponse)
	poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	if err = poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for creation of %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageContainerRead(d, meta)
}

func resourceStorageContainerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageContainerID(d.Id())
	if err != nil {
		return err
	}

	update := blobcontainers.BlobContainer{
		Properties: &blobcontainers.ContainerProperties{},
	}

	if d.HasChange("container_access_type") {
		accessLevelRaw := d.Get("container_access_type").(string)
		update.Properties.PublicAccess = pointer.ToEnum[blobcontainers.PublicAccess](containerAccessTypeConversionMap[accessLevelRaw])
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		update.Properties.Metadata = pointer.To(ExpandMetaData(metaDataRaw))
	}

	if _, err := containerClient.Update(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceStorageContainerRead(d, meta)
}

func resourceStorageContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageContainerID(d.Id())
	if err != nil {
		return err
	}

	existing, err := containerClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			log.Printf("[DEBUG] %q was not found, removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", *id, err)
	}

	if model := existing.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("name", id.ContainerName)
			d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())
			d.Set("container_access_type", containerAccessTypeConversionMap[string(pointer.From(props.PublicAccess))])
			d.Set("default_encryption_scope", props.DefaultEncryptionScope)
			d.Set("encryption_scope_override_enabled", !pointer.From(props.DenyEncryptionScopeOverride))
			d.Set("metadata", FlattenMetaData(pointer.From(props.Metadata)))

			d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
			d.Set("has_legal_hold", props.HasLegalHold)
		}
	}

	account, err := meta.(*clients.Client).Storage.GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Account for Container %q: %v", id, err)
	}

	// Determine the blob endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeBlob)
	if err != nil {
		return fmt.Errorf("determining Blob endpoint: %v", err)
	}

	// Parse the blob endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, meta.(*clients.Client).Storage.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	d.Set("url", containers.NewContainerID(*accountId, id.ContainerName).ID())

	return nil
}

func resourceStorageContainerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageContainerID(d.Id())
	if err != nil {
		return err
	}

	if _, err := containerClient.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %v", d.Id(), err)
	}

	return nil
}
