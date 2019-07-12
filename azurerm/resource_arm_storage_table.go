package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/table/tables"
)

func resourceArmStorageTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageTableCreate,
		Read:   resourceArmStorageTableRead,
		Delete: resourceArmStorageTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		MigrateState:  resourceStorageTableMigrateState,

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

			// TODO: support for ACL's
		},
	}
}

func resourceArmStorageTableCreate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	storageClient := meta.(*ArmClient).storage

	tableName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group: %s", err)
	}

	client, err := storageClient.TablesClient(ctx, *resourceGroup, accountName)
	if err != nil {
		return fmt.Errorf("Error building Table Client: %s", err)
	}

	id := client.GetResourceID(accountName, tableName)
	if requireResourcesToBeImported {
		existing, err := client.Exists(ctx, *resourceGroup, tableName)
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

	d.SetId(id)
	return resourceArmStorageTableRead(d, meta)
}

func resourceArmStorageTableRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).storage
	ctx := meta.(*ArmClient).StopContext

	id, err := tables.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group: %s", err)
	}

	if resourceGroup == nil {
		log.Printf("Unable to determine Resource Group for Storage Account %q (assuming removed)", id.AccountName)
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

	_, err = client.GetACL(ctx, id.AccountName, id.TableName)
	if err != nil {
		return fmt.Errorf("Error retrieving Table %q in Storage Account %q: %s", id.TableName, id.AccountName, err)
	}

	d.Set("name", id.TableName)
	d.Set("storage_account_name", id.AccountName)
	d.Set("resource_group_name", resourceGroup)

	return nil
}

func resourceArmStorageTableDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).storage
	ctx := meta.(*ArmClient).StopContext

	id, err := tables.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group: %s", err)
	}

	if resourceGroup == nil {
		log.Printf("Unable to determine Resource Group for Storage Account %q (assuming removed)", id.AccountName)
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
