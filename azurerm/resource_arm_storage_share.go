package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/schema"
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
		SchemaVersion: 1,
		MigrateState:  resourceStorageShareMigrateState,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageShareName,
			},
			"resource_group_name": resourceGroupNameSchema(),
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5120,
				ValidateFunc: validation.IntBetween(1, 5120),
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceArmStorageShareCreate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", storageAccountName)
	}

	name := d.Get("name").(string)
	metaData := make(map[string]string) // TODO: support MetaData
	options := &storage.FileRequestOptions{}

	log.Printf("[INFO] Creating share %q in storage account %q", name, storageAccountName)
	reference := fileClient.GetShareReference(name)
	err = reference.Create(options)

	log.Printf("[INFO] Setting share %q metadata in storage account %q", name, storageAccountName)
	reference.Metadata = metaData
	reference.SetMetadata(options)

	log.Printf("[INFO] Setting share %q properties in storage account %q", name, storageAccountName)
	reference.Properties = storage.ShareProperties{
		Quota: d.Get("quota").(int),
	}
	reference.SetProperties(options)

	d.SetId(fmt.Sprintf("%s/%s/%s", name, resourceGroupName, storageAccountName))
	return resourceArmStorageShareRead(d, meta)
}

func resourceArmStorageShareRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id := strings.Split(d.Id(), "/")
	name := id[0]
	resourceGroupName := id[1]
	storageAccountName := id[2]

	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing file %q from state", storageAccountName, d.Id())
		d.SetId("")
		return nil
	}

	reference := fileClient.GetShareReference(name)
	exists, err := reference.Exists()
	if err != nil {
		return fmt.Errorf("Error testing existence of share %q: %s", name, err)
	}

	if !exists {
		log.Printf("[INFO] Share %q no longer exists, removing from state...", name)
		d.SetId("")
	}

	url := reference.URL()
	if url == "" {
		log.Printf("[INFO] URL for %q is empty", name)
	}
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	d.Set("storage_account_name", storageAccountName)
	d.Set("url", url)

	reference.FetchAttributes(nil)
	d.Set("quota", reference.Properties.Quota)

	return nil
}

func resourceArmStorageShareUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id := strings.Split(d.Id(), "/")
	name := id[0]
	resourceGroupName := id[1]
	storageAccountName := id[2]

	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", storageAccountName)
	}

	options := &storage.FileRequestOptions{}

	reference := fileClient.GetShareReference(name)

	log.Printf("[INFO] Setting share %q properties in storage account %q", name, storageAccountName)
	reference.Properties = storage.ShareProperties{
		Quota: d.Get("quota").(int),
	}
	reference.SetProperties(options)

	return resourceArmStorageShareRead(d, meta)
}

func resourceArmStorageShareExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id := strings.Split(d.Id(), "/")
	name := id[0]
	resourceGroupName := id[1]
	storageAccountName := id[2]

	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return false, err
	}
	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing share %q from state", storageAccountName, d.Id())
		d.SetId("")
		return false, nil
	}

	log.Printf("[INFO] Checking for existence of share %q.", name)
	reference := fileClient.GetShareReference(name)
	exists, err := reference.Exists()
	if err != nil {
		return false, fmt.Errorf("Error testing existence of share %q: %s", name, err)
	}

	if !exists {
		log.Printf("[INFO] Share %q no longer exists, removing from state...", name)
		d.SetId("")
	}

	return exists, nil
}

func resourceArmStorageShareDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id := strings.Split(d.Id(), "/")
	name := id[0]
	resourceGroupName := id[1]
	storageAccountName := id[2]

	fileClient, accountExists, err := armClient.getFileServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[INFO]Storage Account %q doesn't exist so the file won't exist", storageAccountName)
		return nil
	}

	reference := fileClient.GetShareReference(name)
	options := &storage.FileRequestOptions{}

	if _, err = reference.DeleteIfExists(options); err != nil {
		return fmt.Errorf("Error deleting storage file %q: %s", name, err)
	}

	d.SetId("")
	return nil
}
