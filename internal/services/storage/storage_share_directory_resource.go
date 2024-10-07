// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/file/directories"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/file/shares"
)

func resourceStorageShareDirectory() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceStorageShareDirectoryCreate,
		Read:   resourceStorageShareDirectoryRead,
		Update: resourceStorageShareDirectoryUpdate,
		Delete: resourceStorageShareDirectoryDelete,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := directories.ParseDirectoryID(id, storageDomainSuffix)
			return err
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
				ValidateFunc: validate.StorageShareDirectoryName,
			},

			"storage_share_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageShareDataPlaneID,
			},

			"metadata": MetaDataSchema(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["storage_share_id"].Required = false
		resource.Schema["storage_share_id"].Optional = true
		resource.Schema["storage_share_id"].Computed = true
		resource.Schema["storage_share_id"].ConflictsWith = []string{"share_name", "storage_account_name"}

		resource.Schema["storage_account_name"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			Deprecated:    "the `share_name` and `storage_account_name` properties have been superseded by the `storage_share_id` property and will be removed in version 4.0 of the AzureRM provider",
			ConflictsWith: []string{"storage_share_id"},
			RequiredWith:  []string{"share_name"},
			ValidateFunc:  validation.StringIsNotEmpty,
		}

		resource.Schema["share_name"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			Deprecated:    "the `share_name` and `storage_account_name` properties have been superseded by the `storage_share_id` property and will be removed in version 4.0 of the AzureRM provider",
			ConflictsWith: []string{"storage_share_id"},
			RequiredWith:  []string{"storage_account_name"},
			ValidateFunc:  validation.StringIsNotEmpty,
		}
	}

	return resource
}

func resourceStorageShareDirectoryCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	directoryName := d.Get("name").(string)
	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	var storageShareId *shares.ShareId
	var err error
	if v, ok := d.GetOk("storage_share_id"); ok && v.(string) != "" {
		storageShareId, err = shares.ParseShareID(v.(string), storageClient.StorageDomainSuffix)
		if err != nil {
			return err
		}
	} else if !features.FourPointOhBeta() {
		// TODO: this is needed until `share_name` / `storage_account_name` are removed in favor of `storage_share_id` in v4.0
		// we will retrieve the storage account twice but this will make it easier to refactor later
		storageAccountName := d.Get("storage_account_name").(string)

		account, err := storageClient.FindAccount(ctx, subscriptionId, storageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving Account %q: %v", storageAccountName, err)
		}
		if account == nil {
			return fmt.Errorf("locating Storage Account %q", storageAccountName)
		}

		// Determine the file endpoint, so we can build a data plane ID
		endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeFile)
		if err != nil {
			return fmt.Errorf("determining File endpoint: %v", err)
		}

		// Parse the file endpoint as a data plane account ID
		accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
		if err != nil {
			return fmt.Errorf("parsing Account ID: %v", err)
		}

		storageShareId = pointer.To(shares.NewShareID(*accountId, d.Get("share_name").(string)))
	}

	if storageShareId == nil {
		return fmt.Errorf("determining storage share ID")
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, storageShareId.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Directory %q (Share %q): %v", storageShareId.AccountId.AccountName, directoryName, storageShareId.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", storageShareId.AccountId.AccountName)
	}

	accountId, err := accounts.ParseAccountID(storageShareId.ID(), storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	id := directories.NewDirectoryID(*accountId, storageShareId.ShareName, directoryName)

	client, err := storageClient.FileShareDirectoriesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Directories Client: %v", err)
	}

	existing, err := client.Get(ctx, storageShareId.ShareName, directoryName)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_share_directory", id.ID())
	}

	input := directories.CreateDirectoryInput{
		MetaData: metaData,
	}
	if _, err = client.Create(ctx, storageShareId.ShareName, directoryName, input); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	// Storage Share Directories are eventually consistent
	log.Printf("[DEBUG] Waiting for %s to become available", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   storageShareDirectoryRefreshFunc(ctx, client, id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageShareDirectoryRead(d, meta)
}

func resourceStorageShareDirectoryUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := directories.ParseDirectoryID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Directory %q (Share %q): %v", id.AccountId.AccountName, id.DirectoryPath, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account: %q", id.AccountId.AccountName)
	}

	client, err := storageClient.FileShareDirectoriesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Client: %v", err)
	}

	if _, err = client.SetMetaData(ctx, id.ShareName, id.DirectoryPath, directories.SetMetaDataInput{MetaData: metaData}); err != nil {
		return fmt.Errorf("updating Metadata for %s: %v", id, err)
	}

	return resourceStorageShareDirectoryRead(d, meta)
}

func resourceStorageShareDirectoryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := directories.ParseDirectoryID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Directory %q (Share %q): %v", id.AccountId.AccountName, id.DirectoryPath, id.ShareName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Share Directory %q (Share %s, Account %s) - assuming removed & removing from state", id.DirectoryPath, id.ShareName, id.AccountId.AccountName)
		d.SetId("")
		return nil
	}

	// Determine the file endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeFile)
	if err != nil {
		return fmt.Errorf("determining File endpoint: %v", err)
	}

	// Parse the file endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	storageShareId := shares.NewShareID(*accountId, id.ShareName)

	client, err := storageClient.FileShareDirectoriesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Client: %v", err)
	}

	props, err := client.Get(ctx, id.ShareName, id.DirectoryPath)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("name", id.DirectoryPath)
	d.Set("storage_share_id", storageShareId.ID())

	if !features.FourPointOhBeta() {
		d.Set("storage_account_name", id.AccountId.AccountName)
		d.Set("share_name", id.ShareName)
	}

	if err = d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %v", err)
	}

	return nil
}

func resourceStorageShareDirectoryDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := directories.ParseDirectoryID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Directory %q (Share %q): %v", id.AccountId.AccountName, id.DirectoryPath, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	client, err := storageClient.FileShareDirectoriesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Client: %v", err)
	}

	if _, err = client.Delete(ctx, id.ShareName, id.DirectoryPath); err != nil {
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}

func storageShareDirectoryRefreshFunc(ctx context.Context, client *directories.Client, id directories.DirectoryId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ShareName, id.DirectoryPath)
		if err != nil {
			return nil, strconv.Itoa(res.HttpResponse.StatusCode), fmt.Errorf("retrieving %s: %v", id, err)
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
}
