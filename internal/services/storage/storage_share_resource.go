package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
)

func resourceStorageShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageShareCreate,
		Read:   resourceStorageShareRead,
		Update: resourceStorageShareUpdate,
		Delete: resourceStorageShareDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StorageShareDataPlaneID(id)
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
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	accountName := d.Get("storage_account_name").(string)
	shareName := d.Get("name").(string)
	quota := d.Get("quota").(int)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
	acls := expandStorageShareACLs(aclsRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %s", accountName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Client: %s", err)
	}

	id := parse.NewStorageShareDataPlaneId(accountName, storageClient.Environment.StorageEndpointSuffix, shareName).ID()

	exists, err := client.Exists(ctx, account.ResourceGroup, accountName, shareName)
	if err != nil {
		return fmt.Errorf("checking for existence of existing Storage Share %q (Account %q / Resource Group %q): %+v", shareName, accountName, account.ResourceGroup, err)
	}
	if exists != nil && *exists {
		return tf.ImportAsExistsError("azurerm_storage_share", id)
	}

	log.Printf("[INFO] Creating Share %q in Storage Account %q", shareName, accountName)
	input := shares.CreateInput{
		QuotaInGB:       quota,
		MetaData:        metaData,
		EnabledProtocol: shares.ShareProtocol(d.Get("enabled_protocol").(string)),
	}

	if accessTier := d.Get("access_tier").(string); accessTier != "" {
		tier := shares.AccessTier(accessTier)
		input.AccessTier = &tier
	}

	if err := client.Create(ctx, account.ResourceGroup, accountName, shareName, input); err != nil {
		return fmt.Errorf("creating Share %q (Account %q / Resource Group %q): %+v", shareName, accountName, account.ResourceGroup, err)
	}

	d.SetId(id)
	if err := client.UpdateACLs(ctx, account.ResourceGroup, accountName, shareName, acls); err != nil {
		return fmt.Errorf("setting ACL's for Share %q (Account %q / Resource Group %q): %+v", shareName, accountName, account.ResourceGroup, err)
	}

	return resourceStorageShareRead(d, meta)
}

func resourceStorageShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := parse.StorageShareDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Account %q for Storage Share %q - assuming removed & removing from state", id.AccountName, id.Name)
		d.SetId("")
		return nil
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	props, err := client.Get(ctx, account.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return err
	}
	if props == nil {
		log.Printf("[DEBUG] File Share %q was not found in Account %q / Resource Group %q - assuming removed & removing from state", id.Name, id.AccountName, account.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", id.Name)
	d.Set("storage_account_name", id.AccountName)
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

	resourceManagerId := parse.NewStorageShareResourceManagerID(storageClient.SubscriptionId, account.ResourceGroup, id.AccountName, "default", id.Name)
	d.Set("resource_manager_id", resourceManagerId.ID())

	return nil
}

func resourceStorageShareUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := parse.StorageShareDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	if d.HasChange("quota") {
		log.Printf("[DEBUG] Updating the Quota for File Share %q (Storage Account %q)", id.Name, id.AccountName)
		quota := d.Get("quota").(int)

		if err := client.UpdateQuota(ctx, account.ResourceGroup, id.AccountName, id.Name, quota); err != nil {
			return fmt.Errorf("updating Quota for File Share %q (Storage Account %q): %s", id.Name, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the Quota for File Share %q (Storage Account %q)", id.Name, id.AccountName)
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for File Share %q (Storage Account %q)", id.Name, id.AccountName)

		metaDataRaw := d.Get("metadata").(map[string]interface{})
		metaData := ExpandMetaData(metaDataRaw)

		if err := client.UpdateMetaData(ctx, account.ResourceGroup, id.AccountName, id.Name, metaData); err != nil {
			return fmt.Errorf("updating MetaData for File Share %q (Storage Account %q): %s", id.Name, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the MetaData for File Share %q (Storage Account %q)", id.Name, id.AccountName)
	}

	if d.HasChange("acl") {
		log.Printf("[DEBUG] Updating the ACL's for File Share %q (Storage Account %q)", id.Name, id.AccountName)

		aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
		acls := expandStorageShareACLs(aclsRaw)

		if err := client.UpdateACLs(ctx, account.ResourceGroup, id.AccountName, id.Name, acls); err != nil {
			return fmt.Errorf("updating ACL's for File Share %q (Storage Account %q): %s", id.Name, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the ACL's for File Share %q (Storage Account %q)", id.Name, id.AccountName)
	}

	if d.HasChange("access_tier") {
		log.Printf("[DEBUG] Updating the Access Tier for File Share %q (Storage Account %q)", id.Name, id.AccountName)

		tier := shares.AccessTier(d.Get("access_tier").(string))
		if err := client.UpdateTier(ctx, account.ResourceGroup, id.AccountName, id.Name, tier); err != nil {
			return fmt.Errorf("updating Access Tier for File Share %q (Storage Account %q): %s", id.Name, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the Access Tier for File Share %q (Storage Account %q)", id.Name, id.AccountName)
	}

	return resourceStorageShareRead(d, meta)
}

func resourceStorageShareDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := parse.StorageShareDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Share %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	if err := client.Delete(ctx, account.ResourceGroup, id.AccountName, id.Name); err != nil {
		return fmt.Errorf("deleting File Share %q (Storage Account %q / Resource Group %q): %s", id.Name, id.AccountName, account.ResourceGroup, err)
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
