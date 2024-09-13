// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
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
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/file/shares"
)

func resourceStorageShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageShareCreate,
		Read:   resourceStorageShareRead,
		Update: resourceStorageShareUpdate,
		Delete: resourceStorageShareDelete,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := shares.ParseShareID(id, storageDomainSuffix)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ShareV0ToV1{},
			1: migration.ShareV1ToV2{},
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
				ValidateFunc: validate.StorageShareName,
			},

			"storage_account_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"quota": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 102400),
			},

			"metadata": MetaDataComputedSchema(),

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
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"expiry": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
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

			"enabled_protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(shares.SMB),
					string(shares.NFS),
				}, false),
				Default: string(shares.SMB),
			},

			"resource_manager_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"access_tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(shares.PremiumAccessTier),
						string(shares.HotAccessTier),
						string(shares.CoolAccessTier),
						string(shares.TransactionOptimizedAccessTier),
					}, false),
			},
		},
	}
}

func resourceStorageShareCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountName := d.Get("storage_account_name").(string)
	shareName := d.Get("name").(string)
	quota := d.Get("quota").(int)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
	acls := expandStorageShareACLs(aclsRaw)

	account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %v", accountName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountName)
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

	id := shares.NewShareID(*accountId, shareName)

	protocol := shares.ShareProtocol(d.Get("enabled_protocol").(string))
	if protocol == shares.NFS {
		// Only FileStorage (whose sku tier is Premium only) storage account is able to have NFS file shares.
		// See: https://learn.microsoft.com/en-us/azure/storage/files/storage-files-quick-create-use-linux#applies-to
		if account.Kind != storageaccounts.KindFileStorage {
			return fmt.Errorf("NFS File Share is only supported for Storage Account with kind %q but got `%s`", string(storageaccounts.KindFileStorage), account.Kind)
		}
	}

	// The files API does not support bearer tokens (@manicminer, 2024-02-15)
	client, err := storageClient.FileSharesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
	if err != nil {
		return fmt.Errorf("building File Share Client: %v", err)
	}

	exists, err := client.Exists(ctx, shareName)
	if err != nil {
		return fmt.Errorf("checking for existing %s: %v", id, err)
	}
	if exists != nil && *exists {
		return tf.ImportAsExistsError("azurerm_storage_share", id.ID())
	}

	log.Printf("[INFO] Creating Share %q in Storage Account %q", shareName, accountName)
	input := shares.CreateInput{
		QuotaInGB:       quota,
		MetaData:        metaData,
		EnabledProtocol: protocol,
	}

	if accessTier := d.Get("access_tier").(string); accessTier != "" {
		tier := shares.AccessTier(accessTier)
		input.AccessTier = &tier
	}

	if err = client.Create(ctx, shareName, input); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	if err = client.UpdateACLs(ctx, shareName, shares.SetAclInput{SignedIdentifiers: acls}); err != nil {
		return fmt.Errorf("setting ACLs for %s: %v", id, err)
	}

	return resourceStorageShareRead(d, meta)
}

func resourceStorageShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := shares.ParseShareID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %v", id.AccountId.AccountName, id.ShareName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Account %q for Storage Share %q - assuming removed & removing from state", id.AccountId.AccountName, id.ShareName)
		d.SetId("")
		return nil
	}

	// The files API does not support bearer tokens (@manicminer, 2024-02-15)
	client, err := storageClient.FileSharesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
	if err != nil {
		return fmt.Errorf("building File Share Client for %s: %+v", account.StorageAccountId, err)
	}

	props, err := client.Get(ctx, id.ShareName)
	if err != nil {
		return err
	}
	if props == nil {
		log.Printf("[DEBUG] File Share %q was not found in %s - assuming removed & removing from state", id.ShareName, account.StorageAccountId)
		d.SetId("")
		return nil
	}

	d.Set("name", id.ShareName)
	d.Set("storage_account_name", id.AccountId.AccountName)
	d.Set("quota", props.QuotaGB)
	d.Set("url", id.ID())
	d.Set("enabled_protocol", string(props.EnabledProtocol))

	accessTier := ""
	if props.AccessTier != nil {
		accessTier = string(*props.AccessTier)
	}
	d.Set("access_tier", accessTier)

	if err := d.Set("acl", flattenStorageShareACLs(props.ACLs)); err != nil {
		return fmt.Errorf("flattening `acl`: %+v", err)
	}

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("flattening `metadata`: %+v", err)
	}

	resourceManagerId := parse.NewStorageShareResourceManagerID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, account.StorageAccountId.StorageAccountName, "default", id.ShareName)
	d.Set("resource_manager_id", resourceManagerId.ID())

	return nil
}

func resourceStorageShareUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := shares.ParseShareID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %v", id.AccountId.AccountName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	// The files API does not support bearer tokens (@manicminer, 2024-02-15)
	client, err := storageClient.FileSharesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
	if err != nil {
		return fmt.Errorf("building File Share Client for %s: %+v", account.StorageAccountId, err)
	}

	if d.HasChange("quota") {
		log.Printf("[DEBUG] Updating the Quota for %s", id)
		quota := d.Get("quota").(int)

		if err = client.UpdateQuota(ctx, id.ShareName, quota); err != nil {
			return fmt.Errorf("updating Quota for %s: %v", id, err)
		}

		log.Printf("[DEBUG] Updated the Quota for %s", id)
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for %s", id)

		metaDataRaw := d.Get("metadata").(map[string]interface{})
		metaData := ExpandMetaData(metaDataRaw)

		if err = client.UpdateMetaData(ctx, id.ShareName, metaData); err != nil {
			return fmt.Errorf("updating MetaData for %s: %v", id, err)
		}

		log.Printf("[DEBUG] Updated the MetaData for %s", id)
	}

	if d.HasChange("acl") {
		log.Printf("[DEBUG] Updating the ACLs for %s", id)

		aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
		acls := expandStorageShareACLs(aclsRaw)

		if err = client.UpdateACLs(ctx, id.ShareName, shares.SetAclInput{SignedIdentifiers: acls}); err != nil {
			return fmt.Errorf("updating ACLs for %s: %v", id, err)
		}

		log.Printf("[DEBUG] Updated ACLs for %s", id)
	}

	if d.HasChange("access_tier") {
		log.Printf("[DEBUG] Updating Access Tier for %s", id)

		tier := shares.AccessTier(d.Get("access_tier").(string))
		err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutUpdate), func() *pluginsdk.RetryError {
			err = client.UpdateTier(ctx, id.ShareName, tier)
			if err != nil {
				if strings.Contains(err.Error(), "Cannot change access tier at this moment") {
					return pluginsdk.RetryableError(err)
				}
				return pluginsdk.NonRetryableError(err)
			}
			time.Sleep(30 * time.Second)
			return nil
		})
		if err != nil {
			return fmt.Errorf("updating access tier %s: %+v", id, err)
		}

		log.Printf("[DEBUG] Updated Access Tier for %s", id)
	}

	return resourceStorageShareRead(d, meta)
}

func resourceStorageShareDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := shares.ParseShareID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %v", id.AccountId.AccountName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	// The files API does not support bearer tokens (@manicminer, 2024-02-15)
	client, err := storageClient.FileSharesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
	if err != nil {
		return fmt.Errorf("building File Share Client for %s: %+v", account.StorageAccountId, err)
	}

	if err = client.Delete(ctx, id.ShareName); err != nil {
		if strings.Contains(err.Error(), "The specified share does not exist") {
			return nil
		}
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}

func expandStorageShareACLs(input []interface{}) []shares.SignedIdentifier {
	results := make([]shares.SignedIdentifier, 0)

	for _, v := range input {
		vals := v.(map[string]interface{})

		policies := vals["access_policy"].([]interface{})
		policy := policies[0].(map[string]interface{})

		identifier := shares.SignedIdentifier{
			Id: vals["id"].(string),
			AccessPolicy: shares.AccessPolicy{
				Start:      policy["start"].(string),
				Expiry:     policy["expiry"].(string),
				Permission: policy["permissions"].(string),
			},
		}
		results = append(results, identifier)
	}

	return results
}

func flattenStorageShareACLs(input []shares.SignedIdentifier) []interface{} {
	result := make([]interface{}, 0)

	for _, v := range input {
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
