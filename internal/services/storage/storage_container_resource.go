// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/blobcontainers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
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
	r := &pluginsdk.Resource{
		Create: resourceStorageContainerCreate,
		Read:   resourceStorageContainerRead,
		Delete: resourceStorageContainerDelete,
		Update: resourceStorageContainerUpdate,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			if !features.FivePointOhBeta() {
				if strings.HasPrefix(id, "/subscriptions/") {
					_, err := commonids.ParseStorageContainerID(id)
					return err
				}
				_, err := containers.ParseContainerID(id, storageDomainSuffix)
				return err
			}

			_, err := commonids.ParseStorageContainerID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ContainerV0ToV1{},
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
		},
	}

	if !features.FivePointOhBeta() {
		r.Schema["storage_account_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageAccountName,
			ExactlyOneOf: []string{"storage_account_id", "storage_account_name"},
			Deprecated:   "the `storage_account_name` property has been deprecated in favour of `storage_account_id` and will be removed in version 5.0 of the Provider.",
		}

		r.Schema["storage_account_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			ExactlyOneOf: []string{"storage_account_id", "storage_account_name"},
		}

		r.Schema["resource_manager_id"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Computed:   true,
			Deprecated: "this property has been deprecated in favour of `id` and will be removed in version 5.0 of the Provider.",
		}
	}

	return r
}

func resourceStorageContainerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)
	accessLevelRaw := d.Get("container_access_type").(string)
	accessLevel := expandStorageContainerAccessLevel(accessLevelRaw)
	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	if !features.FivePointOhBeta() {
		storageClient := meta.(*clients.Client).Storage
		if accountName := d.Get("storage_account_name").(string); accountName != "" {
			account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
			if err != nil {
				return fmt.Errorf("retrieving Account %q for Container %q: %v", accountName, containerName, err)
			}
			if account == nil {
				return fmt.Errorf("locating Storage Account %q", accountName)
			}

			containersDataPlaneClient, err := storageClient.ContainersDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building storage client: %v", err)
			}

			// Determine the blob endpoint, so we can build a data plane ID
			endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeBlob)
			if err != nil {
				return fmt.Errorf("determining Blob endpoint: %v", err)
			}

			// Parse the blob endpoint as a data plane account ID
			accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
			if err != nil {
				return fmt.Errorf("parsing Account ID: %v", err)
			}

			id := containers.NewContainerID(*accountId, containerName)

			exists, err := containersDataPlaneClient.Exists(ctx, containerName)
			if err != nil {
				return fmt.Errorf("checking for existing %s: %v", id, err)
			}
			if exists != nil && *exists {
				return tf.ImportAsExistsError("azurerm_storage_container", id.ID())
			}

			log.Printf("[INFO] Creating %s", id)
			input := containers.CreateInput{
				AccessLevel: accessLevel,
				MetaData:    metaData,
			}

			if encryptionScope := d.Get("default_encryption_scope"); encryptionScope.(string) != "" {
				input.DefaultEncryptionScope = encryptionScope.(string)
				input.EncryptionScopeOverrideDisabled = false

				if encryptionScopeOverrideEnabled := d.Get("encryption_scope_override_enabled"); !encryptionScopeOverrideEnabled.(bool) {
					input.EncryptionScopeOverrideDisabled = true
				}
			}

			if err = containersDataPlaneClient.Create(ctx, containerName, input); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}
			d.SetId(id.ID())

			return resourceStorageContainerRead(d, meta)
		}
	}

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := commonids.NewStorageContainerID(subscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, containerName)

	existing, err := containerClient.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %q: %v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_container", id.ID())
	}

	payload := blobcontainers.BlobContainer{
		Properties: &blobcontainers.ContainerProperties{
			PublicAccess: pointer.To(blobcontainers.PublicAccess(containerAccessTypeConversionMap[accessLevelRaw])),
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

	if _, err = containerClient.Create(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageContainerRead(d, meta)
}

func resourceStorageContainerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !features.FivePointOhBeta() && !strings.HasPrefix(d.Id(), "/subscriptions/") {
		storageClient := meta.(*clients.Client).Storage
		id, err := containers.ParseContainerID(d.Id(), storageClient.StorageDomainSuffix)
		if err != nil {
			return err
		}

		account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
		if err != nil {
			return fmt.Errorf("retrieving Account %q for Container %q: %v", id.AccountId.AccountName, id.ContainerName, err)
		}
		if account == nil {
			return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
		}
		if d.HasChange("container_access_type") {
			log.Printf("[DEBUG] Updating Access Level for %s...", id)

			// Updating metadata does not work with AAD authentication, returns a cryptic 404
			client, err := storageClient.ContainersDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
			if err != nil {
				return fmt.Errorf("building Containers Client: %v", err)
			}

			accessLevelRaw := d.Get("container_access_type").(string)
			accessLevel := expandStorageContainerAccessLevel(accessLevelRaw)

			if err = client.UpdateAccessLevel(ctx, id.ContainerName, accessLevel); err != nil {
				return fmt.Errorf("updating Access Level for %s: %v", id, err)
			}

			log.Printf("[DEBUG] Updated Access Level for %s", id)
		}

		if d.HasChange("metadata") {
			log.Printf("[DEBUG] Updating Metadata for %s...", id)

			client, err := storageClient.ContainersDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Containers Client: %v", err)
			}

			metaDataRaw := d.Get("metadata").(map[string]interface{})
			metaData := ExpandMetaData(metaDataRaw)

			if err = client.UpdateMetaData(ctx, id.ContainerName, metaData); err != nil {
				return fmt.Errorf("updating Metadata for %s: %v", id, err)
			}

			log.Printf("[DEBUG] Updated Metadata for %s", id)
		}

		return resourceStorageContainerRead(d, meta)
	}

	id, err := commonids.ParseStorageContainerID(d.Id())
	if err != nil {
		return err
	}

	update := blobcontainers.BlobContainer{
		Properties: &blobcontainers.ContainerProperties{},
	}

	if d.HasChange("container_access_type") {
		accessLevelRaw := d.Get("container_access_type").(string)
		update.Properties.PublicAccess = pointer.To(blobcontainers.PublicAccess(containerAccessTypeConversionMap[accessLevelRaw]))
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !features.FivePointOhBeta() && !strings.HasPrefix(d.Id(), "/subscriptions/") {
		storageClient := meta.(*clients.Client).Storage
		id, err := containers.ParseContainerID(d.Id(), storageClient.StorageDomainSuffix)
		if err != nil {
			return err
		}

		account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
		if err != nil {
			return fmt.Errorf("retrieving Account %q for Container %q: %v", id.AccountId.AccountName, id.ContainerName, err)
		}
		if account == nil {
			log.Printf("[DEBUG] Unable to locate Account %q for Storage Container %q - assuming removed & removing from state", id.AccountId.AccountName, id.ContainerName)
			d.SetId("")
			return nil
		}

		client, err := storageClient.ContainersDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Containers Client: %v", err)
		}

		props, err := client.Get(ctx, id.ContainerName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %v", id, err)
		}
		if props == nil {
			log.Printf("[DEBUG] Container %q was not found in %s - assuming removed & removing from state", id.ContainerName, id.AccountId)
			d.SetId("")
			return nil
		}

		d.Set("name", id.ContainerName)
		d.Set("storage_account_name", id.AccountId.AccountName)

		d.Set("container_access_type", flattenStorageContainerAccessLevel(props.AccessLevel))

		d.Set("default_encryption_scope", props.DefaultEncryptionScope)
		d.Set("encryption_scope_override_enabled", !props.EncryptionScopeOverrideDisabled)

		if err = d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
			return fmt.Errorf("setting `metadata`: %v", err)
		}

		d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
		d.Set("has_legal_hold", props.HasLegalHold)

		resourceManagerId := commonids.NewStorageContainerID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, id.AccountId.AccountName, id.ContainerName)
		d.Set("resource_manager_id", resourceManagerId.ID())

		return nil
	}

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
			if !features.FivePointOhBeta() {
				d.Set("resource_manager_id", id.ID())
			}
		}
	}

	return nil
}

func resourceStorageContainerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !features.FivePointOhBeta() && !strings.HasPrefix(d.Id(), "/subscriptions/") {
		storageClient := meta.(*clients.Client).Storage

		id, err := containers.ParseContainerID(d.Id(), storageClient.StorageDomainSuffix)
		if err != nil {
			return err
		}

		account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
		if err != nil {
			return fmt.Errorf("retrieving Account %q for Container %q: %v", id.AccountId.AccountName, id.ContainerName, err)
		}
		if account == nil {
			return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
		}

		client, err := storageClient.ContainersDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Containers Client: %v", err)
		}

		if err = client.Delete(ctx, id.ContainerName); err != nil {
			return fmt.Errorf("deleting %s: %v", id, err)
		}

		return nil
	}

	id, err := commonids.ParseStorageContainerID(d.Id())
	if err != nil {
		return err
	}

	if _, err := containerClient.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %v", d.Id(), err)
	}

	return nil
}

func expandStorageContainerAccessLevel(input string) containers.AccessLevel {
	// for historical reasons, "private" above is an empty string in the API
	// so the enum doesn't 1:1 match. You could argue the SDK should handle this
	// but this is suitable for now
	if input == "private" {
		return containers.Private
	}

	return containers.AccessLevel(input)
}

func flattenStorageContainerAccessLevel(input containers.AccessLevel) string {
	// for historical reasons, "private" above is an empty string in the API
	if input == containers.Private {
		return "private"
	}

	return string(input)
}
