package azurerm

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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
			"resource_group_name": resourceGroupNameSchema(),
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
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

func resourceArmStorageTableCreate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext
	environment := armClient.environment

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	tableClient, accountExists, err := armClient.getTableServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", storageAccountName)
	}

	table := tableClient.GetTableReference(name)
	id := fmt.Sprintf("https://%s.table.%s/%s", storageAccountName, environment.StorageEndpointSuffix, name)

	if requireResourcesToBeImported {
		metaDataLevel := storage.MinimalMetadata
		options := &storage.QueryTablesOptions{}
		tables, e := tableClient.QueryTables(metaDataLevel, options)
		if e != nil {
			return fmt.Errorf("Error checking if Table %q exists (Account %q / Resource Group %q): %s", name, storageAccountName, resourceGroupName, e)
		}

		for _, table := range tables.Tables {
			if table.Name == name {
				return tf.ImportAsExistsError("azurerm_storage_table", id)
			}
		}
	}

	log.Printf("[INFO] Creating table %q in storage account %q.", name, storageAccountName)
	timeout := uint(60)
	options := &storage.TableOptions{}
	err = table.Create(timeout, storage.NoMetadata, options)
	if err != nil {
		return fmt.Errorf("Error creating table %q in storage account %q: %s", name, storageAccountName, err)
	}

	d.SetId(id)
	return resourceArmStorageTableRead(d, meta)
}

func resourceArmStorageTableRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := parseStorageTableID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	if err != nil {
		return err
	}

	if resourceGroup == nil {
		log.Printf("Unable to determine Resource Group for Storage Account %q (assuming removed)", id.storageAccountName)
		d.SetId("")
		return nil
	}

	tableClient, accountExists, err := armClient.getTableServiceClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return err
	}

	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing table %q from state", id.storageAccountName, id.tableName)
		d.SetId("")
		return nil
	}

	metaDataLevel := storage.MinimalMetadata
	options := &storage.QueryTablesOptions{}
	tables, err := tableClient.QueryTables(metaDataLevel, options)
	if err != nil {
		return fmt.Errorf("Failed to retrieve Tables in Storage Account %q: %s", id.tableName, err)
	}

	var storageTable *storage.Table
	for _, table := range tables.Tables {
		if table.Name == id.tableName {
			storageTable = &table
			break
		}
	}

	if storageTable == nil {
		log.Printf("[INFO] Table %q does not exist in Storage Account %q, removing from state...", id.tableName, id.storageAccountName)
		d.SetId("")
		return nil
	}

	d.Set("name", id.tableName)
	d.Set("storage_account_name", id.storageAccountName)
	d.Set("resource_group_name", resourceGroup)

	return nil
}

func resourceArmStorageTableDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := parseStorageTableID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	if err != nil {
		return err
	}

	if resourceGroup == nil {
		log.Printf("Unable to determine Resource Group for Storage Account %q (assuming removed)", id.storageAccountName)
		return nil
	}

	tableClient, accountExists, err := armClient.getTableServiceClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[INFO] Storage Account %q doesn't exist so the table won't exist", id.storageAccountName)
		return nil
	}

	table := tableClient.GetTableReference(id.tableName)
	timeout := uint(60)
	options := &storage.TableOptions{}

	log.Printf("[INFO] Deleting Table %q in Storage Account %q", id.tableName, id.storageAccountName)
	if err := table.Delete(timeout, options); err != nil {
		return fmt.Errorf("Error deleting table %q from Storage Account %q: %s", id.tableName, id.storageAccountName, err)
	}

	return nil
}

type storageTableId struct {
	storageAccountName string
	tableName          string
}

func parseStorageTableID(input string) (*storageTableId, error) {
	// https://myaccount.table.core.windows.net/table1
	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("Error parsing %q as a URI: %+v", input, err)
	}

	segments := strings.Split(uri.Host, ".")
	if len(segments) > 0 {
		storageAccountName := segments[0]
		table := strings.Replace(uri.Path, "/", "", 1)
		id := storageTableId{
			storageAccountName: storageAccountName,
			tableName:          table,
		}
		return &id, nil
	}

	return nil, nil
}
