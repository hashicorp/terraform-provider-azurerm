// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
)

var storageAccountStaticWebSitePropertiesResourceName = "azurerm_storage_account_static_website_properties"

func resourceStorageAccountStaticWebSiteProperties() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountStaticWebSitePropertiesCreate,
		Read:   resourceStorageAccountStaticWebSitePropertiesRead,
		Update: resourceStorageAccountStaticWebSitePropertiesUpdate,
		Delete: resourceStorageAccountStaticWebSitePropertiesDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseStorageAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"error_404_document": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"index_document": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceStorageAccountStaticWebSitePropertiesCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountStaticWebSitePropertiesResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountStaticWebSitePropertiesResourceName)

	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	// NOTE: Import error cannot be supported for this resource...

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Kind == nil {
		return fmt.Errorf("retrieving %s: `model.Kind` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}

	if existing.Model.Sku == nil {
		return fmt.Errorf("retrieving %s: `model.Sku` was nil", id)
	}

	dataPlaneAccount, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if dataPlaneAccount == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	var accountKind storageaccounts.Kind
	var accountTier storageaccounts.SkuTier
	accountReplicationType := ""

	accountKind = *existing.Model.Kind
	accountReplicationType = strings.Split(string(existing.Model.Sku.Name), "_")[1]
	if existing.Model.Sku.Tier != nil {
		accountTier = *existing.Model.Sku.Tier
	}

	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)
	if err := waitForDataPlaneToBecomeAvailableForAccount(ctx, storageClient, dataPlaneAccount, supportLevel); err != nil {
		return fmt.Errorf("waiting for the Data Plane for %s to become available: %+v", id, err)
	}

	if val, ok := d.GetOk("properties"); ok {
		if !supportLevel.supportStaticWebsite {
			return fmt.Errorf("static websites are not supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *dataPlaneAccount, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Accounts Data Plane Client: %s", err)
		}

		staticWebsiteProps := expandAccountStaticWebsiteProperties(val.([]interface{}))

		if _, err = accountsClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceStorageAccountStaticWebSitePropertiesRead(d, meta)
}

func resourceStorageAccountStaticWebSitePropertiesUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountStaticWebSitePropertiesResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountStaticWebSitePropertiesResourceName)

	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Kind == nil {
		return fmt.Errorf("retrieving %s: `model.Kind` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}

	if existing.Model.Sku == nil {
		return fmt.Errorf("retrieving %s: `model.Sku` was nil", id)
	}

	var accountKind storageaccounts.Kind
	var accountTier storageaccounts.SkuTier
	accountReplicationType := ""

	accountKind = *existing.Model.Kind
	accountReplicationType = strings.Split(string(existing.Model.Sku.Name), "_")[1]
	if existing.Model.Sku.Tier != nil {
		accountTier = *existing.Model.Sku.Tier
	}

	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)

	if d.HasChange("properties") {
		if !supportLevel.supportStaticWebsite {
			return fmt.Errorf("static website properties are not supported for a storage account with the account kind %q in sku tier %q", accountKind, accountTier)
		}

		account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		if account == nil {
			return fmt.Errorf("unable to locate %s", *id)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Data Plane client for %s: %+v", *id, err)
		}

		staticWebsiteProps := expandAccountStaticWebsiteProperties(d.Get("properties").([]interface{}))

		if _, err = accountsClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
			return fmt.Errorf("updating Static Website Properties for %s: %+v", *id, err)
		}
	}

	return resourceStorageAccountStaticWebSitePropertiesRead(d, meta)
}

func resourceStorageAccountStaticWebSitePropertiesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// we then need to find the storage account
	account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Accounts Data Plane Client: %s", err)
	}

	staticWebsiteProps, err := accountsClient.GetServiceProperties(ctx, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving static website properties for %s: %+v", *id, err)
	}

	staticWebsiteProperties := flattenAccountStaticWebsiteProperties(staticWebsiteProps)

	if err := d.Set("properties", staticWebsiteProperties); err != nil {
		return fmt.Errorf("setting `properties`: %+v", err)
	}

	return nil
}

func resourceStorageAccountStaticWebSitePropertiesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountStaticWebSitePropertiesResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountStaticWebSitePropertiesResourceName)

	_, err = client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// NOTE: Since this is a fake resource the best we can do
	// is clear the settings and set enabled to false...
	staticWebsiteProps := accounts.StorageServiceProperties{
		StaticWebsite: &accounts.StaticWebsite{
			Enabled: false,
		},
	}

	dataPlaneAccount, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if dataPlaneAccount == nil {
		return fmt.Errorf("unable to locate %s", *id)
	}

	accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *dataPlaneAccount, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Plane Accounts Client %s: %+v", *id, err)
	}

	if _, err = accountsClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAccountStaticWebsiteProperties(input []interface{}) accounts.StorageServiceProperties {
	properties := accounts.StorageServiceProperties{
		StaticWebsite: &accounts.StaticWebsite{
			Enabled: false,
		},
	}
	if len(input) == 0 {
		return properties
	}

	properties.StaticWebsite.Enabled = true

	// @tombuildsstuff: this looks weird, doesn't it?
	// Since the presence of this block signifies the website's enabled however all fields within it are optional
	// TF Core returns a nil object when there's no keys defined within the block, rather than an empty map. As
	// such this hack allows us to have a Static Website block with only Enabled configured, without the optional
	// inner properties.
	if val := input[0]; val != nil {
		attr := val.(map[string]interface{})
		if v, ok := attr["index_document"]; ok {
			properties.StaticWebsite.IndexDocument = v.(string)
		}

		if v, ok := attr["error_404_document"]; ok {
			properties.StaticWebsite.ErrorDocument404Path = v.(string)
		}
	}

	return properties
}

func flattenAccountStaticWebsiteProperties(input accounts.GetServicePropertiesResult) []interface{} {
	if staticWebsite := input.StaticWebsite; staticWebsite != nil {
		if !staticWebsite.Enabled {
			return []interface{}{}
		}

		return []interface{}{
			map[string]interface{}{
				"error_404_document": staticWebsite.ErrorDocument404Path,
				"index_document":     staticWebsite.IndexDocument,
			},
		}
	}
	return []interface{}{}
}
