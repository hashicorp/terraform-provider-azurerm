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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/fileshares"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageShareRm() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageShareRmCreate,
		Read:   resourceStorageShareRmRead,
		Update: resourceStorageShareRmUpdate,
		Delete: resourceStorageShareRmDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := fileshares.ParseShareID(id)
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
				ValidateFunc: validate.StorageShareName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

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
									"start_time": {
										Type:             pluginsdk.TypeString,
										Optional:         true,
										DiffSuppressFunc: suppress.RFC3339Time,
										ValidateFunc:     validation.IsRFC3339Time,
									},
									"expiry_time": {
										Type:             pluginsdk.TypeString,
										Optional:         true,
										DiffSuppressFunc: suppress.RFC3339Time,
										ValidateFunc:     validation.IsRFC3339Time,
									},
									"permission": {
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
					string(fileshares.EnabledProtocolsSMB),
					string(fileshares.EnabledProtocolsNFS),
				}, false),
				Default: string(fileshares.EnabledProtocolsSMB),
			},

			"root_squash": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(fileshares.RootSquashTypeAllSquash),
					string(fileshares.RootSquashTypeNoRootSquash),
					string(fileshares.RootSquashTypeRootSquash),
				}, false),
			},

			"access_tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(fileshares.ShareAccessTierPremium),
						string(fileshares.ShareAccessTierHot),
						string(fileshares.ShareAccessTierCool),
						string(fileshares.ShareAccessTierTransactionOptimized),
					}, false),
			},
		},
	}
}

func resourceStorageShareRmCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	fileSharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	accountName := d.Get("storage_account_name").(string)
	shareName := d.Get("name").(string)
	quota := d.Get("quota").(int)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
	acls := expandStorageShareRmACLs(aclsRaw)

	account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %v", accountName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountName)
	}

	shareId := fileshares.NewShareID(subscriptionId, resourceGroupName, accountName, shareName)
	existing, err := fileSharesClient.Get(ctx, shareId, fileshares.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("could not get fileshare")

		}
	}

	log.Printf("[INFO] Creating Share %q in Storage Account %q", shareName, accountName)

	fileShare := fileshares.FileShare{
		Name: pointer.To(shareName),
		Properties: &fileshares.FileShareProperties{
			SignedIdentifiers: pointer.To(acls),
			Metadata:          &metaData,
			ShareQuota:        pointer.To(int64(quota)),
		},
	}

	if protocol := d.Get("enabled_protocol").(string); protocol != "" {
		enabledProtocol := parseProtocol(protocol)
		if enabledProtocol == fileshares.EnabledProtocolsNFS {
			// Only FileStorage (whose sku tier is Premium only) storage account is able to have NFS file shares.
			// See: https://learn.microsoft.com/en-us/azure/storage/files/storage-files-quick-create-use-linux#applies-to
			if account.Kind != storageaccounts.KindFileStorage {
				return fmt.Errorf("NFS File Share is only supported for Storage Account with kind %q but got `%s`", string(storageaccounts.KindFileStorage), account.Kind)
			}

			if rootSquash := d.Get("root_squash").(string); rootSquash != "" {
				fileShare.Properties.RootSquash = pointer.To(parseRootSquashType(rootSquash))
			}

			if quota < 100 {
				return fmt.Errorf("NFS File Share required a quota of 100 at least but got %d", quota)
			}
		}

		fileShare.Properties.EnabledProtocols = pointer.To(enabledProtocol)
	}

	if accessTier := d.Get("access_tier").(string); accessTier != "" {
		fileShare.Properties.AccessTier = pointer.To(parseAccessTier(accessTier))
	}

	createResp, err := fileSharesClient.Create(ctx, shareId, fileShare, fileshares.DefaultCreateOperationOptions())
	if err != nil {
		return fmt.Errorf("creating %s: %v", shareId.ID(), err)
	}

	id, err := fileshares.ParseShareID(pointer.From(createResp.Model.Id))
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceStorageShareRmRead(d, meta)
}

func resourceStorageShareRmRead(d *pluginsdk.ResourceData, meta interface{}) error {
	fileSharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fileshares.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	resp, err := fileSharesClient.Get(ctx, *id, fileshares.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ShareName)
	d.Set("storage_account_name", id.StorageAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	d.Set("quota", int(pointer.From(resp.Model.Properties.ShareQuota)))

	if resp.Model.Properties.EnabledProtocols != nil {
		d.Set("enabled_protocol", string(*resp.Model.Properties.EnabledProtocols))
	}

	accessTier := ""
	if resp.Model.Properties.AccessTier != nil {
		accessTier = string(pointer.From(resp.Model.Properties.AccessTier))
	}
	d.Set("access_tier", accessTier)

	if resp.Model.Properties.SignedIdentifiers != nil {
		if err := d.Set("acl", flattenStorageShareRmACLs(pointer.From(resp.Model.Properties.SignedIdentifiers))); err != nil {
			return fmt.Errorf("flattening `acl`: %+v", err)
		}
	}

	metadata := make(map[string]interface{})
	if resp.Model.Properties.Metadata != nil {
		metadata = FlattenMetaData(pointer.From(resp.Model.Properties.Metadata))
	}

	d.Set("root_squash", string(pointer.From(resp.Model.Properties.RootSquash)))

	d.Set("metadata", metadata)

	return nil
}

func resourceStorageShareRmUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	fileSharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fileshares.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	fileShare := fileshares.FileShare{
		Properties: &fileshares.FileShareProperties{},
	}

	if d.HasChange("quota") {
		log.Printf("[DEBUG] Updating the Quota for %s", id)
		quota := int64(d.Get("quota").(int))

		fileShare.Properties.ShareQuota = pointer.To(quota)
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for %s", id)

		metaDataRaw := d.Get("metadata").(map[string]interface{})
		metaData := ExpandMetaData(metaDataRaw)

		fileShare.Properties.Metadata = &metaData
	}

	if d.HasChange("acl") {
		log.Printf("[DEBUG] Updating the ACLs for %s", id)

		aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
		acls := expandStorageShareRmACLs(aclsRaw)

		fileShare.Properties.SignedIdentifiers = pointer.To(acls)
	}

	if d.HasChange("access_tier") {
		log.Printf("[DEBUG] Updating Access Tier for %s", id)

		accessTier := d.Get("access_tier").(string)
		fileShare.Properties.AccessTier = pointer.To(parseAccessTier(accessTier))
	}

	if d.HasChange("root_squash") {
		log.Printf("[DEBUG] Updating Root Squash for %s", id)

		rootSquash := d.Get("root_squash").(string)
		fileShare.Properties.RootSquash = pointer.To(parseRootSquashType(rootSquash))
	}

	if _, err = fileSharesClient.Update(ctx, *id, fileShare); err != nil {
		return fmt.Errorf("could not update storage share %v", err)
	}

	return resourceStorageShareRmRead(d, meta)
}

func resourceStorageShareRmDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	fileSharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fileshares.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	if _, err = fileSharesClient.Delete(ctx, *id, fileshares.DefaultDeleteOperationOptions()); err != nil {
		if strings.Contains(err.Error(), "The specified share does not exist") {
			return nil
		}
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}

func parseProtocol(input string) fileshares.EnabledProtocols {
	vals := map[string]fileshares.EnabledProtocols{
		"nfs": fileshares.EnabledProtocolsNFS,
		"smb": fileshares.EnabledProtocolsSMB,
	}

	return vals[strings.ToLower(input)]
}

func parseAccessTier(input string) fileshares.ShareAccessTier {
	vals := map[string]fileshares.ShareAccessTier{
		"cool":                 fileshares.ShareAccessTierCool,
		"hot":                  fileshares.ShareAccessTierHot,
		"premium":              fileshares.ShareAccessTierPremium,
		"transactionoptimized": fileshares.ShareAccessTierTransactionOptimized,
	}

	return vals[strings.ToLower(input)]
}

func parseRootSquashType(input string) fileshares.RootSquashType {
	vals := map[string]fileshares.RootSquashType{
		"allsquash":    fileshares.RootSquashTypeAllSquash,
		"norootsquash": fileshares.RootSquashTypeNoRootSquash,
		"rootsquash":   fileshares.RootSquashTypeRootSquash,
	}

	return vals[strings.ToLower(input)]
}

func expandStorageShareRmACLs(input []interface{}) []fileshares.SignedIdentifier {
	results := make([]fileshares.SignedIdentifier, 0)

	for _, v := range input {
		vals := v.(map[string]interface{})

		policies := vals["access_policy"].([]interface{})
		policy := policies[0].(map[string]interface{})

		identifier := fileshares.SignedIdentifier{
			Id: pointer.To(vals["id"].(string)),
			AccessPolicy: &fileshares.AccessPolicy{
				StartTime:  pointer.To(policy["start_time"].(string)),
				ExpiryTime: pointer.To(policy["expiry_time"].(string)),
				Permission: pointer.To(policy["permission"].(string)),
			},
		}
		results = append(results, identifier)
	}

	return results
}

func flattenStorageShareRmACLs(input []fileshares.SignedIdentifier) []interface{} {
	result := make([]interface{}, 0)

	for _, v := range input {
		output := map[string]interface{}{
			"id": v.Id,
			"access_policy": []interface{}{
				map[string]interface{}{
					"start_time":  v.AccessPolicy.StartTime,
					"expiry_time": v.AccessPolicy.ExpiryTime,
					"permission":  v.AccessPolicy.Permission,
				},
			},
		}

		result = append(result, output)
	}

	return result
}
