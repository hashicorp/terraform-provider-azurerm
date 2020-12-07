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
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
)

func resourceArmStorageTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageTableCreate,
		Read:   resourceArmStorageTableRead,
		Delete: resourceArmStorageTableDelete,
		Update: resourceArmStorageTableUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 2,
		StateUpgraders: []schema.StateUpgrader{
			{
				// this should have been applied from pre-0.12 migration system; backporting just in-case
				Type:    ResourceStorageTableStateResourceV0V1().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceStorageTableStateUpgradeV0ToV1,
				Version: 0,
			},
			{
				Type:    ResourceStorageTableStateResourceV0V1().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceStorageTableStateUpgradeV1ToV2,
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
				ValidateFunc: validate.StorageTableName,
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateArmStorageAccountName,
			},

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
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"expiry": {
										Type:         schema.TypeString,
										Required:     true,
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
		},
	}
}

func resourceArmStorageTableCreate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	tableName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)
	aclsRaw := d.Get("acl").(*schema.Set).List()
	acls := expandStorageTableACLs(aclsRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %s", accountName, tableName, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q!", accountName)
	}

	client, err := storageClient.TablesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Table Client: %s", err)
	}

	id := parse.NewStorageTableDataPlaneId(accountName, storageClient.Environment.StorageEndpointSuffix, tableName).ID()

	exists, err := client.Exists(ctx, account.ResourceGroup, accountName, tableName)
	if err != nil {
		return fmt.Errorf("checking for existence of existing Storage Table %q (Account %q / Resource Group %q): %+v", tableName, accountName, account.ResourceGroup, err)
	}
	if exists != nil && *exists {
		return tf.ImportAsExistsError("azurerm_storage_table", id)
	}

	log.Printf("[DEBUG] Creating Table %q in Storage Account %q.", tableName, accountName)
	if err := client.Create(ctx, account.ResourceGroup, accountName, tableName); err != nil {
		return fmt.Errorf("creating Table %q within Storage Account %q: %s", tableName, accountName, err)
	}

	d.SetId(id)
	if err := client.UpdateACLs(ctx, account.ResourceGroup, accountName, tableName, acls); err != nil {
		return fmt.Errorf("setting ACL's for Storage Table %q (Account %q / Resource Group %q): %+v", tableName, accountName, account.ResourceGroup, err)
	}

	return resourceArmStorageTableRead(d, meta)
}

func resourceArmStorageTableRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageTableDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Table %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		log.Printf("Unable to determine Resource Group for Storage Storage Table %q (Account %s) - assuming removed & removing from state", id.Name, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TablesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Table Client: %s", err)
	}

	exists, err := client.Exists(ctx, account.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Table %q (Storage Account %q / Resource Group %q): %s", id.Name, id.AccountName, account.ResourceGroup, err)
	}
	if exists == nil || !*exists {
		log.Printf("[DEBUG] Storage Account %q not found, removing table %q from state", id.AccountName, id.Name)
		d.SetId("")
		return nil
	}

	acls, err := client.GetACLs(ctx, account.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving ACL's %q in Storage Account %q: %s", id.Name, id.AccountName, err)
	}

	d.Set("name", id.Name)
	d.Set("storage_account_name", id.AccountName)

	if err := d.Set("acl", flattenStorageTableACLs(acls)); err != nil {
		return fmt.Errorf("flattening `acl`: %+v", err)
	}

	return nil
}

func resourceArmStorageTableDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageTableDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Table %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.TablesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Table Client: %s", err)
	}

	log.Printf("[INFO] Deleting Table %q in Storage Account %q", id.Name, id.AccountName)
	if err := client.Delete(ctx, account.ResourceGroup, id.AccountName, id.Name); err != nil {
		return fmt.Errorf("deleting Table %q from Storage Account %q: %s", id.Name, id.AccountName, err)
	}

	return nil
}

func resourceArmStorageTableUpdate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageTableDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.TablesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Table Client: %s", err)
	}

	if d.HasChange("acl") {
		log.Printf("[DEBUG] Updating the ACL's for Storage Table %q (Storage Account %q)", id.Name, id.AccountName)

		aclsRaw := d.Get("acl").(*schema.Set).List()
		acls := expandStorageTableACLs(aclsRaw)

		if err := client.UpdateACLs(ctx, account.ResourceGroup, id.AccountName, id.Name, acls); err != nil {
			return fmt.Errorf("updating ACL's for Table %q (Storage Account %q): %s", id.Name, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the ACL's for Storage Table %q (Storage Account %q)", id.Name, id.AccountName)
	}

	return resourceArmStorageTableRead(d, meta)
}

func expandStorageTableACLs(input []interface{}) []tables.SignedIdentifier {
	results := make([]tables.SignedIdentifier, 0)

	for _, v := range input {
		vals := v.(map[string]interface{})

		policies := vals["access_policy"].([]interface{})
		policy := policies[0].(map[string]interface{})

		identifier := tables.SignedIdentifier{
			Id: vals["id"].(string),
			AccessPolicy: tables.AccessPolicy{
				Start:      policy["start"].(string),
				Expiry:     policy["expiry"].(string),
				Permission: policy["permissions"].(string),
			},
		}
		results = append(results, identifier)
	}

	return results
}

func flattenStorageTableACLs(input *[]tables.SignedIdentifier) []interface{} {
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
