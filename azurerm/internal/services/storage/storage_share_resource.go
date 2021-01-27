package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/shares"
)

func resourceStorageShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageShareCreate,
		Read:   resourceStorageShareRead,
		Update: resourceStorageShareUpdate,
		Delete: resourceStorageShareDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 2,
		StateUpgraders: []schema.StateUpgrader{
			{
				// this should have been applied from pre-0.12 migration system; backporting just in-case
				Type:    resourceStorageShareStateResourceV0V1().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceStorageShareStateUpgradeV0ToV1,
				Version: 0,
			},
			{
				Type:    resourceStorageShareStateResourceV0V1().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceStorageShareStateUpgradeV1ToV2,
				Version: 1,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageShareName,
			},

			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"quota": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5120,
				ValidateFunc: validation.IntBetween(1, 102400),
			},

			"metadata": MetaDataComputedSchema(),

			"acl": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},
						"access_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"expiry": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"permissions": {
										Type:         schema.TypeString,
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
				Type:     schema.TypeString,
				Computed: true,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageShareCreate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	accountName := d.Get("storage_account_name").(string)
	shareName := d.Get("name").(string)
	quota := d.Get("quota").(int)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	aclsRaw := d.Get("acl").(*schema.Set).List()
	acls := expandStorageShareACLs(aclsRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Share %q: %s", accountName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Client: %s", err)
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
		QuotaInGB: quota,
		MetaData:  metaData,
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

func resourceStorageShareRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceStorageShareUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := parse.StorageShareDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Share %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
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
			return fmt.Errorf("Error updating MetaData for File Share %q (Storage Account %q): %s", id.Name, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the MetaData for File Share %q (Storage Account %q)", id.Name, id.AccountName)
	}

	if d.HasChange("acl") {
		log.Printf("[DEBUG] Updating the ACL's for File Share %q (Storage Account %q)", id.Name, id.AccountName)

		aclsRaw := d.Get("acl").(*schema.Set).List()
		acls := expandStorageShareACLs(aclsRaw)

		if err := client.UpdateACLs(ctx, account.ResourceGroup, id.AccountName, id.Name, acls); err != nil {
			return fmt.Errorf("Error updating ACL's for File Share %q (Storage Account %q): %s", id.Name, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the ACL's for File Share %q (Storage Account %q)", id.Name, id.AccountName)
	}

	return resourceStorageShareRead(d, meta)
}

func resourceStorageShareDelete(d *schema.ResourceData, meta interface{}) error {
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
