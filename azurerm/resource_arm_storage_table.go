package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/table/tables"
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
				Type:    resourceStorageTableStateResourceV0V1().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceStorageTableStateUpgradeV0ToV1,
				Version: 0,
			},
			{
				Type:    resourceStorageTableStateResourceV0V1().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceStorageTableStateUpgradeV1ToV2,
				Version: 1,
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageTableName,
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageAccountName,
			},

			// TODO: deprecate this in the docs
			"resource_group_name": azure.SchemaResourceGroupNameDeprecated(),

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
										ValidateFunc: validate.NoEmptyStrings,
									},
									"expiry": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
									"permissions": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
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
	ctx := meta.(*ArmClient).StopContext
	storageClient := meta.(*ArmClient).Storage

	tableName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)
	aclsRaw := d.Get("acl").(*schema.Set).List()
	acls := expandStorageTableACLs(aclsRaw)

	resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Table %q (Account %s): %s", tableName, accountName, err)
	}
	if resourceGroup == nil {
		return fmt.Errorf("Unable to locate Resource Group for Storage Share %q (Account %s) - assuming removed & removing from state", tableName, accountName)
	}

	client, err := storageClient.TablesClient(ctx, *resourceGroup, accountName)
	if err != nil {
		return fmt.Errorf("Error building Table Client: %s", err)
	}

	id := client.GetResourceID(accountName, tableName)
	if features.ShouldResourcesBeImported() {
		existing, err := client.Exists(ctx, accountName, tableName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing) {
				return fmt.Errorf("Error checking for existence of existing Storage Table %q (Account %q / Resource Group %q): %+v", tableName, accountName, *resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing) {
			return tf.ImportAsExistsError("azurerm_storage_table", id)
		}
	}

	log.Printf("[DEBUG] Creating Table %q in Storage Account %q.", tableName, accountName)
	if _, err := client.Create(ctx, accountName, tableName); err != nil {
		return fmt.Errorf("Error creating Table %q within Storage Account %q: %s", tableName, accountName, err)
	}

	if _, err := client.SetACL(ctx, accountName, tableName, acls); err != nil {
		return fmt.Errorf("Error setting ACL's for Storage Table %q (Account %q / Resource Group %q): %+v", tableName, accountName, *resourceGroup, err)
	}

	d.SetId(id)
	return resourceArmStorageTableRead(d, meta)
}

func resourceArmStorageTableRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).Storage
	ctx := meta.(*ArmClient).StopContext

	id, err := tables.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Table %q (Account %s): %s", id.TableName, id.AccountName, err)
	}

	if resourceGroup == nil {
		log.Printf("Unable to determine Resource Group for Storage Storage Table %q (Account %s) - assuming removed & removing from state", id.TableName, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TablesClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Table Client: %s", err)
	}

	exists, err := client.Exists(ctx, id.AccountName, id.TableName)
	if err != nil {
		if utils.ResponseWasNotFound(exists) {
			log.Printf("[DEBUG] Storage Account %q not found, removing table %q from state", id.AccountName, id.TableName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Table %q in Storage Account %q: %s", id.TableName, id.AccountName, err)
	}

	acls, err := client.GetACL(ctx, id.AccountName, id.TableName)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL's %q in Storage Account %q: %s", id.TableName, id.AccountName, err)
	}

	d.Set("name", id.TableName)
	d.Set("storage_account_name", id.AccountName)
	d.Set("resource_group_name", resourceGroup)

	if err := d.Set("acl", flattenStorageTableACLs(acls)); err != nil {
		return fmt.Errorf("Error flattening `acl`: %+v", err)
	}

	return nil
}

func resourceArmStorageTableDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).Storage
	ctx := meta.(*ArmClient).StopContext

	id, err := tables.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Table %q (Account %s): %s", id.TableName, id.AccountName, err)
	}

	if resourceGroup == nil {
		log.Printf("Unable to determine Resource Group for Storage Storage Table %q (Account %s) - assuming removed & removing from state", id.TableName, id.AccountName)
		return nil
	}

	client, err := storageClient.TablesClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Table Client: %s", err)
	}

	log.Printf("[INFO] Deleting Table %q in Storage Account %q", id.TableName, id.AccountName)
	if _, err := client.Delete(ctx, id.AccountName, id.TableName); err != nil {
		return fmt.Errorf("Error deleting Table %q from Storage Account %q: %s", id.TableName, id.AccountName, err)
	}

	return nil
}

func resourceArmStorageTableUpdate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).Storage
	ctx := meta.(*ArmClient).StopContext

	id, err := tables.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Table %q (Account %s): %s", id.TableName, id.AccountName, err)
	}

	if resourceGroup == nil {
		log.Printf("Unable to determine Resource Group for Storage Storage Table %q (Account %s) - assuming removed & removing from state", id.TableName, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TablesClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Table Client: %s", err)
	}

	if d.HasChange("acl") {
		log.Printf("[DEBUG] Updating the ACL's for Storage Table %q (Storage Account %q)", id.TableName, id.AccountName)

		aclsRaw := d.Get("acl").(*schema.Set).List()
		acls := expandStorageTableACLs(aclsRaw)

		if _, err := client.SetACL(ctx, id.AccountName, id.TableName, acls); err != nil {
			return fmt.Errorf("Error updating ACL's for Storage Table %q (Storage Account %q): %s", id.TableName, id.AccountName, err)
		}

		log.Printf("[DEBUG] Updated the ACL's for Storage Table %q (Storage Account %q)", id.TableName, id.AccountName)
	}

	return resourceArmStorageTableRead(d, meta)
}

func validateArmStorageTableName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if value == "table" {
		errors = append(errors, fmt.Errorf(
			"Table Storage %q cannot use the word `table`: %q",
			k, value))
	}
	if !regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]{2,62}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"Table Storage %q cannot begin with a numeric character, only alphanumeric characters are allowed and must be between 3 and 63 characters long: %q",
			k, value))
	}

	return warnings, errors
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

func flattenStorageTableACLs(input tables.GetACLResult) []interface{} {
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
