package storage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/shares"
)

func resourceArmStorageShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageShareCreate,
		Read:   resourceArmStorageShareRead,
		Update: resourceArmStorageShareUpdate,
		Delete: resourceArmStorageShareDelete,
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
				ValidateFunc: validate.ShareName,
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
func resourceArmStorageShareCreate(d *schema.ResourceData, meta interface{}) error {
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

	id := client.GetResourceID(accountName, shareName)

	existing, err := client.GetProperties(ctx, accountName, shareName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for existence of existing Storage Share %q (Account %q / Resource Group %q): %+v", shareName, accountName, account.ResourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_storage_share", id)
	}

	log.Printf("[INFO] Creating Share %q in Storage Account %q", shareName, accountName)
	input := shares.CreateInput{
		QuotaInGB: quota,
		MetaData:  metaData,
	}

	if resp, err := client.Create(ctx, accountName, shareName, input); err != nil {
		// If we fail due to previous delete still in progress, then we can retry
		if utils.ResponseWasConflict(resp) && strings.Contains(err.Error(), "ShareBeingDeleted") {
			stateConf := &resource.StateChangeConf{
				Pending:        []string{"waitingOnDelete"},
				Target:         []string{"succeeded"},
				Refresh:        storageShareCreateRefreshFunc(ctx, client, accountName, shareName, input),
				PollInterval:   10 * time.Second,
				NotFoundChecks: 180,
				Timeout:        d.Timeout(schema.TimeoutCreate),
			}

			if _, err := stateConf.WaitForState(); err != nil {
				return fmt.Errorf("failed creating share: %+v", err)
			}
		} else {
			return fmt.Errorf("failed creating share: %+v", err)
		}
	}

	if _, err := client.SetACL(ctx, accountName, shareName, acls); err != nil {
		return fmt.Errorf("Error setting ACL's for Share %q (Account %q / Resource Group %q): %+v", shareName, accountName, account.ResourceGroup, err)
	}

	d.SetId(id)
	return resourceArmStorageShareRead(d, meta)
}

func resourceArmStorageShareRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := shares.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Share %q: %s", id.AccountName, id.ShareName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Account %q for Storage Share %q - assuming removed & removing from state", id.AccountName, id.ShareName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	props, err := client.GetProperties(ctx, id.AccountName, id.ShareName)
	if err != nil {
		if utils.ResponseWasNotFound(props.Response) {
			log.Printf("[DEBUG] File Share %q was not found in Account %q / Resource Group %q - assuming removed & removing from state", id.ShareName, id.AccountName, account.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving File Share %q (Account %q / Resource Group %q): %s", id.ShareName, id.AccountName, account.ResourceGroup, err)
	}

	acls, err := client.GetACL(ctx, id.AccountName, id.ShareName)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL's for File Share %q (Account %q / Resource Group %q): %s", id.ShareName, id.AccountName, account.ResourceGroup, err)
	}

	d.Set("name", id.ShareName)
	d.Set("storage_account_name", id.AccountName)
	d.Set("url", client.GetResourceID(id.AccountName, id.ShareName))
	d.Set("quota", props.ShareQuota)

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("Error flattening `metadata`: %+v", err)
	}

	if err := d.Set("acl", flattenStorageShareACLs(acls)); err != nil {
		return fmt.Errorf("Error flattening `acl`: %+v", err)
	}

	resourceManagerId := client.GetResourceManagerResourceID(storageClient.SubscriptionId, account.ResourceGroup, id.AccountName, id.ShareName)
	d.Set("resource_manager_id", resourceManagerId)

	return nil
}

func resourceArmStorageShareUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := shares.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Share %q: %s", id.AccountName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	if d.HasChange("quota") {
		log.Printf("[DEBUG] Updating the Quota for File Share %q (Storage Account %q)", id.ShareName, id.AccountName)
		quota := d.Get("quota").(int)
		if _, err := client.SetProperties(ctx, id.AccountName, id.ShareName, quota); err != nil {
			return fmt.Errorf("Error updating Quota for File Share %q (Storage Account %q): %s", id.ShareName, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the Quota for File Share %q (Storage Account %q)", id.ShareName, id.AccountName)
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for File Share %q (Storage Account %q)", id.ShareName, id.AccountName)

		metaDataRaw := d.Get("metadata").(map[string]interface{})
		metaData := ExpandMetaData(metaDataRaw)

		if _, err := client.SetMetaData(ctx, id.AccountName, id.ShareName, metaData); err != nil {
			return fmt.Errorf("Error updating MetaData for File Share %q (Storage Account %q): %s", id.ShareName, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the MetaData for File Share %q (Storage Account %q)", id.ShareName, id.AccountName)
	}

	if d.HasChange("acl") {
		log.Printf("[DEBUG] Updating the ACL's for File Share %q (Storage Account %q)", id.ShareName, id.AccountName)

		aclsRaw := d.Get("acl").(*schema.Set).List()
		acls := expandStorageShareACLs(aclsRaw)

		if _, err := client.SetACL(ctx, id.AccountName, id.ShareName, acls); err != nil {
			return fmt.Errorf("Error updating ACL's for File Share %q (Storage Account %q): %s", id.ShareName, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the ACL's for File Share %q (Storage Account %q)", id.ShareName, id.AccountName)
	}

	return resourceArmStorageShareRead(d, meta)
}

func resourceArmStorageShareDelete(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := shares.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Share %q: %s", id.AccountName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building File Share Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	deleteSnapshots := true
	if _, err := client.Delete(ctx, id.AccountName, id.ShareName, deleteSnapshots); err != nil {
		return fmt.Errorf("Error deleting File Share %q (Storage Account %q / Resource Group %q): %s", id.ShareName, id.AccountName, account.ResourceGroup, err)
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

func flattenStorageShareACLs(input shares.GetACLResult) []interface{} {
	result := make([]interface{}, 0)

	for _, v := range input.SignedIdentifiers {
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

func storageShareCreateRefreshFunc(ctx context.Context, client *shares.Client, accountName string, shareName string, input shares.CreateInput) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Create(ctx, accountName, shareName, input)
		if err != nil {
			if !utils.ResponseWasConflict(resp) {
				return nil, "", err
			}

			if utils.ResponseWasConflict(resp) && strings.Contains(err.Error(), "ShareBeingDeleted") {
				return nil, "waitingOnDelete", nil
			}
		}

		return "succeeded", "succeeded", nil
	}
}
