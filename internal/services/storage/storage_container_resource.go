// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/containers"
)

func resourceStorageContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageContainerCreate,
		Read:   resourceStorageContainerRead,
		Delete: resourceStorageContainerDelete,
		Update: resourceStorageContainerUpdate,

		// Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
		// 	_, err := parse.StorageContainerDataPlaneID(id)
		// 	return err
		// }),

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

			"storage_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountName,
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

			"metadata": MetaDataComputedSchema(),

			// TODO: support for ACL's, Legal Holds and Immutability Policies
			"has_immutability_policy": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"has_legal_hold": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"resource_manager_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageContainerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	domainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("retrieving domain suffix for Storage Accounts")
	}

	containerName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)
	accessLevelRaw := d.Get("container_access_type").(string)
	accessLevel := expandStorageContainerAccessLevel(accessLevelRaw)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Container %q: %s", accountName, containerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building storage client: %+v", err)
	}

	accountId, err := accounts.ParseAccountID(account.ID, *domainSuffix)
	if err != nil {
		return fmt.Errorf("parsing account ID %s: %+v", account.ID, err)
	}

	id := containers.NewContainerID(*accountId, containerName).ID()
	exists, err := client.Exists(ctx, account.ResourceGroup, containerName)
	if err != nil {
		return err
	}
	if exists != nil && *exists {
		return tf.ImportAsExistsError("azurerm_storage_container", id)
	}

	log.Printf("[INFO] Creating Container %q in Storage Account %q", containerName, accountName)
	input := containers.CreateInput{
		AccessLevel: accessLevel,
		MetaData:    metaData,
	}

	if err := client.Create(ctx, account.ResourceGroup, containerName, input); err != nil {
		return fmt.Errorf("failed creating container: %+v", err)
	}

	d.SetId(id)
	return resourceStorageContainerRead(d, meta)
}

func resourceStorageContainerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	domainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("retrieving domain suffix for Storage Accounts")
	}

	id, err := containers.ParseContainerID(d.Id(), *domainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Container %q: %s", id.AccountId.AccountName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountId.AccountName)
	}
	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Containers Client for Storage Account %q (Resource Group %q): %s", id.AccountId.AccountName, account.ResourceGroup, err)
	}

	if d.HasChange("container_access_type") {
		log.Printf("[DEBUG] Updating the Access Control for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountId.AccountName, account.ResourceGroup)
		accessLevelRaw := d.Get("container_access_type").(string)
		accessLevel := expandStorageContainerAccessLevel(accessLevelRaw)

		if err := client.UpdateAccessLevel(ctx, account.ResourceGroup, id.ContainerName, accessLevel); err != nil {
			return fmt.Errorf("updating the Access Control for Container %q (Storage Account %q / Resource Group %q): %s", id.ContainerName, id.AccountId.AccountName, account.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Updated the Access Control for Container %q (Storage Account %q / Resource Group %q)", id.ContainerName, id.AccountId.AccountName, account.ResourceGroup)
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountId.AccountName, account.ResourceGroup)
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		metaData := ExpandMetaData(metaDataRaw)

		if err := client.UpdateMetaData(ctx, account.ResourceGroup, id.ContainerName, metaData); err != nil {
			return fmt.Errorf("updating the MetaData for Container %q (Storage Account %q / Resource Group %q): %s", id.ContainerName, id.AccountId.AccountName, account.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Updated the MetaData for Container %q (Storage Account %q / Resource Group %q)", id.ContainerName, id.AccountId.AccountName, account.ResourceGroup)
	}

	return resourceStorageContainerRead(d, meta)
}

func resourceStorageContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	domainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("retrieving domain suffix for Storage Accounts")
	}

	id, err := containers.ParseContainerID(d.Id(), *domainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Container %q: %s", id.AccountId.AccountName, id.ContainerName, err)
	}
	if account == nil {
		log.Printf("[DEBUG] Unable to locate Account %q for Storage Container %q - assuming removed & removing from state", id.AccountId.AccountName, id.ContainerName)
		d.SetId("")
		return nil
	}
	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Containers Client for Storage Account %q (Resource Group %q): %s", id.AccountId.AccountName, account.ResourceGroup, err)
	}

	props, err := client.Get(ctx, account.ResourceGroup, id.ContainerName)
	if err != nil {
		return fmt.Errorf("retrieving Container %q (Account %q / Resource Group %q): %s", id.ContainerName, id.AccountId.AccountName, account.ResourceGroup, err)
	}
	if props == nil {
		log.Printf("[DEBUG] Container %q was not found in Account %q / Resource Group %q - assuming removed & removing from state", id.ContainerName, id.AccountId.AccountName, account.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", id.ContainerName)
	d.Set("storage_account_name", id.AccountId.AccountName)

	d.Set("container_access_type", flattenStorageContainerAccessLevel(props.AccessLevel))

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %+v", err)
	}

	d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
	d.Set("has_legal_hold", props.HasLegalHold)

	resourceManagerId := commonids.NewStorageContainerID(subscriptionId, account.ResourceGroup, id.AccountId.AccountName, id.ContainerName)
	d.Set("resource_manager_id", resourceManagerId.ID())

	return nil
}

func resourceStorageContainerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	domainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("retrieving domain suffix for Storage Accounts")
	}

	id, err := containers.ParseContainerID(d.Id(), *domainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Container %q: %s", id.AccountId.AccountName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountId.AccountName)
	}
	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Containers Client for Storage Account %q (Resource Group %q): %s", id.AccountId.AccountName, account.ResourceGroup, err)
	}

	if err := client.Delete(ctx, account.ResourceGroup, id.ContainerName); err != nil {
		return fmt.Errorf("deleting Container %q (Storage Account %q / Resource Group %q): %s", id.ContainerName, id.AccountId.AccountName, account.ResourceGroup, err)
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
