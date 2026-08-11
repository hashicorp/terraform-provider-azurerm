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
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

func resourceStorageShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageShareCreate,
		Read:   resourceStorageShareRead,
		Update: resourceStorageShareUpdate,
		Delete: resourceStorageShareDelete,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := fileshares.ParseShareID(id)
			return err
		}),

		SchemaVersion: 3,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ShareV0ToV1{},
			1: migration.ShareV1ToV2{},
			2: migration.StorageShareV2ToV3{},
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

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"quota": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 102400),
			},

			"metadata": MetaDataComputedSchema(),

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
										ValidateFunc: validation.IsRFC3339Time,
									},
									"expiry": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsRFC3339Time,
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
					}, false,
				),
			},

			"rbac_scope_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageShareCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	sharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := fileshares.NewShareID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := sharesClient.Get(ctx, id, fileshares.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_share", id.ID())
		}
	}

	payload := fileshares.FileShare{
		Properties: &fileshares.FileShareProperties{
			EnabledProtocols:  pointer.To(fileshares.EnabledProtocols(d.Get("enabled_protocol").(string))),
			Metadata:          pointer.To(ExpandMetaData(d.Get("metadata").(map[string]interface{}))),
			ShareQuota:        pointer.To(int64(d.Get("quota").(int))),
			SignedIdentifiers: expandStorageShareACLs(d.Get("acl").(*pluginsdk.Set).List()),
		},
	}

	if sharedAccessTier, ok := d.GetOk("access_tier"); ok && sharedAccessTier.(string) != "" {
		payload.Properties.AccessTier = pointer.To(fileshares.ShareAccessTier(sharedAccessTier.(string)))
	}

	pollerType := custompollers.NewStorageShareCreatePoller(sharesClient, id, payload)
	poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	if err = poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageShareRead(d, meta)
}

func resourceStorageShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	sharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fileshares.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	existing, err := sharesClient.Get(ctx, *id, fileshares.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			log.Printf("[DEBUG] %q was not found, removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", *id, err)
	}

	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())
	d.Set("name", id.ShareName)

	if model := existing.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("quota", props.ShareQuota)
			// Resource Manager treats nil and "SMB" as the same and we may not get a full response here
			enabledProtocols := fileshares.EnabledProtocolsSMB
			if props.EnabledProtocols != nil {
				enabledProtocols = *props.EnabledProtocols
			}
			d.Set("enabled_protocol", string(enabledProtocols))
			d.Set("access_tier", string(pointer.From(props.AccessTier)))
			d.Set("acl", flattenStorageShareACLs(pointer.From(props.SignedIdentifiers)))
			d.Set("metadata", FlattenMetaData(pointer.From(props.Metadata)))
		}
	}

	// TODO - The following section for `url` will need to be updated to go-azure-sdk when the Giovanni Deprecation process has been completed
	account, err := meta.(*clients.Client).Storage.GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Account for Share %q: %v", id, err)
	}

	// Determine the file endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeFile)
	if err != nil {
		return fmt.Errorf("determining File endpoint: %v", err)
	}

	// Parse the file endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, meta.(*clients.Client).Storage.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	d.Set("url", shares.NewShareID(*accountId, id.ShareName).ID())
	d.Set("rbac_scope_id", parse.NewStorageShareResourceManagerID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, "default", id.ShareName).ID())

	return nil
}

func resourceStorageShareUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	sharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fileshares.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	update := fileshares.FileShare{
		Properties: &fileshares.FileShareProperties{},
	}

	if d.HasChange("quota") {
		quota := d.Get("quota").(int)
		update.Properties.ShareQuota = pointer.To(int64(quota))
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		metaData := ExpandMetaData(metaDataRaw)

		update.Properties.Metadata = pointer.To(metaData)
	}

	if d.HasChange("acl") {
		update.Properties.SignedIdentifiers = expandStorageShareACLs(d.Get("acl").(*pluginsdk.Set).List())
	}

	if d.HasChange("access_tier") {
		tier := shares.AccessTier(d.Get("access_tier").(string))
		update.Properties.AccessTier = pointer.To(fileshares.ShareAccessTier(tier))
	}

	if _, err = sharesClient.Update(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceStorageShareRead(d, meta)
}

func resourceStorageShareDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	fileSharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fileshares.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := fileSharesClient.Delete(ctx, *id, fileshares.DefaultDeleteOperationOptions()); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %v", id, err)
		}
	}

	return nil
}

func expandStorageShareACLs(input []interface{}) *[]fileshares.SignedIdentifier {
	results := make([]fileshares.SignedIdentifier, 0)

	for _, v := range input {
		acl := v.(map[string]interface{})

		policies := acl["access_policy"].([]interface{})
		policy := policies[0].(map[string]interface{})

		identifier := fileshares.SignedIdentifier{
			Id: pointer.To(acl["id"].(string)),
			AccessPolicy: &fileshares.AccessPolicy{
				StartTime:  pointer.To(policy["start"].(string)),
				ExpiryTime: pointer.To(policy["expiry"].(string)),
				Permission: pointer.To(policy["permissions"].(string)),
			},
		}
		results = append(results, identifier)
	}

	return pointer.To(results)
}

func flattenStorageShareACLs(input []fileshares.SignedIdentifier) []interface{} {
	result := make([]interface{}, 0)

	for _, v := range input {
		output := map[string]interface{}{
			"id": v.Id,
			"access_policy": []interface{}{
				map[string]interface{}{
					"start":       v.AccessPolicy.StartTime,
					"expiry":      v.AccessPolicy.ExpiryTime,
					"permissions": v.AccessPolicy.Permission,
				},
			},
		}

		result = append(result, output)
	}

	return result
}
