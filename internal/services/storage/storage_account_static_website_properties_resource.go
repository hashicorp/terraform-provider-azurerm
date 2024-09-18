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
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
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

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:CREATE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	// TODO: Add Import error supported for this resource...

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	accountTier := pointer.From(model.Sku.Tier)
	accountKind := pointer.From(model.Kind)
	replicationType := strings.ToUpper(strings.Split(string(model.Sku.Name), "_")[1])
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, replicationType)

	if !supportLevel.supportStaticWebsite {
		return fmt.Errorf("%q are not supported for account kind %q", storageAccountStaticWebSitePropertiesResourceName, accountKind)
	}

	// NOTE: Wait for the static website data plane container to become available...
	log.Printf("[DEBUG] [%s:CREATE] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	dataPlaneAccount, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if dataPlaneAccount == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	// NOTE: Wait for the data plane static website container to become available...
	log.Printf("[DEBUG] [%s:CREATE] Calling 'custompollers.NewDataPlaneStaticWebsiteAvailabilityPoller' building Static Website Poller: %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	pollerType, err := custompollers.NewDataPlaneStaticWebsiteAvailabilityPoller(ctx, storageClient, dataPlaneAccount)
	if err != nil {
		return fmt.Errorf("waiting for the Data Plane for %s to become available: building Static Website Poller: %+v", id, err)
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'pollers.NewPoller: %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	log.Printf("[DEBUG] [%s:CREATE] Calling 'poller.PollUntilDone' building Accounts Data Plane Client: %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for the Data Plane for %s to become available: waiting for the Static Website to become available: %+v", id, err)
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'storageClient.AccountsDataPlaneClient' building Accounts Data Plane Client: %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	accountsDataPlaneClient, err := storageClient.AccountsDataPlaneClient(ctx, *dataPlaneAccount, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Accounts Data Plane Client: %s", err)
	}

	// Wrap the flattened schema into an interface slice to reuse the same expand/flatten functions...
	properties := make(map[string]interface{})

	if v, ok := d.GetOk("error_404_document"); ok {
		properties["error_404_document"] = v.(string)
	}

	if v, ok := d.GetOk("index_document"); ok {
		properties["index_document"] = v.(string)
	}

	staticWebsiteProps := expandAccountStaticWebsiteProperties([]interface{}{properties})

	log.Printf("[DEBUG] [%s:CREATE] Calling 'accountsDataPlaneClient.SetServiceProperties': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	if _, err = accountsDataPlaneClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
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

	log.Printf("[DEBUG] [%s:UPDATE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	accountTier := pointer.From(model.Sku.Tier)
	accountKind := pointer.From(model.Kind)
	replicationType := strings.ToUpper(strings.Split(string(model.Sku.Name), "_")[1])
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, replicationType)

	if d.HasChange("error_404_document") || d.HasChange("index_document") {
		if !supportLevel.supportStaticWebsite {
			return fmt.Errorf("%q are not supported for account kind %q", storageAccountStaticWebSitePropertiesResourceName, accountKind)
		}

		log.Printf("[DEBUG] [%s:UPDATE] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
		account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		if account == nil {
			return fmt.Errorf("unable to locate %s", *id)
		}

		log.Printf("[DEBUG] [%s:UPDATE] Calling 'storageClient.AccountsDataPlaneClient': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
		accountsDataPlaneClient, err := storageClient.AccountsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Data Plane client for %s: %+v", *id, err)
		}

		// Wrap the flattened schema into an interface slice to reuse the existing legacy expand/flatten functions...
		staticWebsiteProperties := make(map[string]interface{}, 0)

		if v, ok := d.GetOk("error_404_document"); ok {
			staticWebsiteProperties["error_404_document"] = v.(string)
		}

		if v, ok := d.GetOk("index_document"); ok {
			staticWebsiteProperties["index_document"] = v.(string)
		}

		staticWebsiteProps := expandAccountStaticWebsiteProperties([]interface{}{staticWebsiteProperties})

		log.Printf("[DEBUG] [%s:UPDATE] Calling 'accountsDataPlaneClient.SetServiceProperties': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
		if _, err = accountsDataPlaneClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
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

	log.Printf("[DEBUG] [%s:READ] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// we then need to find the storage account
	log.Printf("[DEBUG] [%s:READ] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	log.Printf("[DEBUG] [%s:READ] Calling 'storageClient.AccountsDataPlaneClient': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	accountsDataPlaneClient, err := storageClient.AccountsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Accounts Data Plane Client: %s", err)
	}

	log.Printf("[DEBUG] [%s:READ] Calling 'accountsDataPlaneClient.GetServiceProperties': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	staticWebsiteProps, err := accountsDataPlaneClient.GetServiceProperties(ctx, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving static website properties for %s: %+v", *id, err)
	}

	staticWebsiteProperties := flattenAccountStaticWebsiteProperties(staticWebsiteProps)

	var indexDocument string
	var error404Document string

	// Pull the values out of the flattened properties slice to set the individual flattened values to state...
	if val := staticWebsiteProperties[0]; val != nil {
		attr := val.(map[string]interface{})
		if v, ok := attr["index_document"]; ok {
			indexDocument = v.(string)
		}

		if v, ok := attr["error_404_document"]; ok {
			error404Document = v.(string)
		}
	}

	if err := d.Set("index_document", indexDocument); err != nil {
		return fmt.Errorf("index_document `properties`: %+v", err)
	}

	if err := d.Set("error_404_document", error404Document); err != nil {
		return fmt.Errorf("setting `error_404_document`: %+v", err)
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

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:DELETE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	_, err = client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// NOTE: Since this is a fake resource that has been split off from the main storage account resource
	// the best we can do is reset the values to the default settings...
	staticWebsiteProps := defaultStaticWebsiteProperties()

	log.Printf("[DEBUG] [%s:DELETE] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	dataPlaneAccount, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if dataPlaneAccount == nil {
		return fmt.Errorf("unable to locate %s", *id)
	}

	log.Printf("[DEBUG] [%s:DELETE] Calling 'storageClient.AccountsDataPlaneClient': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *dataPlaneAccount, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Plane Accounts Client %s: %+v", *id, err)
	}

	log.Printf("[DEBUG] [%s:DELETE] Calling 'accountsClient.SetServiceProperties': %s", strings.ToUpper(storageAccountStaticWebSitePropertiesResourceName), id)
	if _, err = accountsClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAccountStaticWebsiteProperties(input []interface{}) accounts.StorageServiceProperties {
	properties := defaultStaticWebsiteProperties()

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

func defaultStaticWebsiteProperties() accounts.StorageServiceProperties {
	return accounts.StorageServiceProperties{
		StaticWebsite: &accounts.StaticWebsite{
			Enabled: false,
		},
	}
}
