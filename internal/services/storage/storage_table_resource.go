// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	legacyTables "github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func resourceStorageTable() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageTableCreate,
		Read:   resourceStorageTableRead,
		Delete: resourceStorageTableDelete,
		Update: resourceStorageTableUpdate,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := tables.ParseTableID(id)
			return err
		}),

		SchemaVersion: 3,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.TableV0ToV1{},
			1: migration.TableV1ToV2{},
			2: migration.TableV2ToV3{},
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
				ValidateFunc: validate.StorageTableName,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"acl": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},
						"access_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"start": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"expiry": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"permissions": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},

			"resource_manager_id": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "The Resource Manager ID of this Storage Table.",
			},
		},
	}
}

func resourceStorageTableCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tableName := d.Get("name").(string)
	aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
	acls := expandStorageTableACLs(aclsRaw)

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}
	accountName := accountId.StorageAccountName

	account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %s", accountName, tableName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountName)
	}

	tablesDataPlaneClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %s", err)
	}

	id := parse.NewStorageTableResourceManagerID(subscriptionId, accountId.ResourceGroupName, accountName, "default", tableName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		exists, err := tablesDataPlaneClient.Exists(ctx, tableName)
		if err != nil {
			return fmt.Errorf("checking for existing %s: %v", id, err)
		}
		if exists != nil && *exists {
			return tf.ImportAsExistsError("azurerm_storage_table", id.ID())
		}
	}

	if err = tablesDataPlaneClient.Create(ctx, tableName); err != nil {
		return fmt.Errorf("creating %s: %v", id.ID(), err)
	}

	d.SetId(id.ID())

	aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %v", err)
	}

	if err = aclClient.UpdateACLs(ctx, tableName, acls); err != nil {
		return fmt.Errorf("setting ACLs for %s: %v", id.ID(), err)
	}

	return resourceStorageTableRead(d, meta)
}

func resourceStorageTableRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rmId, err := parse.StorageTableResourceManagerID(d.Id())
	if err != nil {
		return err
	}

	tableName := rmId.TableName
	accountName := rmId.StorageAccountName

	account, err := storageClient.GetAccount(ctx, commonids.NewStorageAccountID(rmId.SubscriptionId, rmId.ResourceGroup, rmId.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", accountName, tableName, err)
	}
	if account == nil {
		log.Printf("Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", tableName, accountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %v", err)
	}

	exists, err := client.Exists(ctx, tableName)
	if err != nil {
		return fmt.Errorf("retrieving table %q: %v", tableName, err)
	}
	if exists == nil || !*exists {
		log.Printf("[DEBUG] table %q not found, removing from state", tableName)
		d.SetId("")
		return nil
	}

	aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %v", err)
	}

	acls, err := aclClient.GetACLs(ctx, tableName)
	if err != nil {
		return fmt.Errorf("retrieving ACLs for table %q: %v", tableName, err)
	}

	d.Set("name", tableName)
	d.Set("storage_account_id", commonids.NewStorageAccountID(rmId.SubscriptionId, rmId.ResourceGroup, rmId.StorageAccountName).ID())

	if err = d.Set("acl", flattenStorageTableACLs(acls)); err != nil {
		return fmt.Errorf("setting `acl`: %v", err)
	}

	return nil
}

func resourceStorageTableDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rmId, err := parse.StorageTableResourceManagerID(d.Id())
	if err != nil {
		return err
	}

	var account *client.AccountDetails
	if meta.(*clients.Client).Storage.StorageUseAzureAD {
		account = &client.AccountDetails{
			StorageAccountId: commonids.NewStorageAccountID(rmId.SubscriptionId, rmId.ResourceGroup, rmId.StorageAccountName),
		}
	} else {
		account, err = storageClient.GetAccount(ctx, commonids.NewStorageAccountID(rmId.SubscriptionId, rmId.ResourceGroup, rmId.StorageAccountName))
		if err != nil {
			return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", rmId.StorageAccountName, rmId.TableName, err)
		}
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", rmId.StorageAccountName)
	}

	client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %v", err)
	}

	if err = client.Delete(ctx, rmId.TableName); err != nil {
		if strings.Contains(err.Error(), "unexpected status 40") {
			return nil
		}
		return fmt.Errorf("deleting table %q: %v", rmId.TableName, err)
	}

	return nil
}

func resourceStorageTableUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rmId, err := tables.ParseTableID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, rmId.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", rmId.StorageAccountName, rmId.TableName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", rmId.StorageAccountName)
	}

	if d.HasChange("acl") {
		aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
		acls := expandStorageTableACLs(aclsRaw)

		aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Tables Client: %v", err)
		}

		if err = aclClient.UpdateACLs(ctx, rmId.TableName, acls); err != nil {
			return fmt.Errorf("updating ACLs for table %q: %v", rmId.TableName, err)
		}
	}

	return resourceStorageTableRead(d, meta)
}

func expandStorageTableACLs(input []interface{}) []legacyTables.SignedIdentifier {
	results := make([]legacyTables.SignedIdentifier, 0)

	for _, v := range input {
		vals := v.(map[string]interface{})

		policies := vals["access_policy"].([]interface{})
		policy := policies[0].(map[string]interface{})

		identifier := legacyTables.SignedIdentifier{
			Id: vals["id"].(string),
			AccessPolicy: legacyTables.AccessPolicy{
				Start:      policy["start"].(string),
				Expiry:     policy["expiry"].(string),
				Permission: policy["permissions"].(string),
			},
		}
		results = append(results, identifier)
	}

	return results
}

func flattenStorageTableACLs(input *[]legacyTables.SignedIdentifier) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		output := map[string]interface{}{
			"id": v.Id,
			"access_policy": []interface{}{
				map[string]interface{}{
					"start":       v.AccessPolicy.Start,
					"expiry":      v.AccessPolicy.Expiry,
					"permissions": v.AccessPolicy.Permission,
				},
			},
		}

		result = append(result, output)
	}

	return result
}
