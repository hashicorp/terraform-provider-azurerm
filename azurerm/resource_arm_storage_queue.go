package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/queue/queues"
)

func resourceArmStorageQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageQueueCreate,
		Read:   resourceArmStorageQueueRead,
		Update: resourceArmStorageQueueUpdate,
		Delete: resourceArmStorageQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		MigrateState:  resourceStorageQueueMigrateState,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageQueueName,
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageAccountName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDeprecated(),

			"metadata": storage.MetaDataSchema(),
		},
	}
}

func validateArmStorageQueueName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q", k))
	}

	if regexp.MustCompile(`^-`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q cannot start with a hyphen", k))
	}

	if regexp.MustCompile(`-$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q cannot end with a hyphen", k))
	}

	if len(value) > 63 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 63 characters", k))
	}

	if len(value) < 3 {
		errors = append(errors, fmt.Errorf(
			"%q must be at least 3 characters", k))
	}

	return warnings, errors
}

func resourceArmStorageQueueCreate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).storage
	ctx := meta.(*ArmClient).StopContext

	queueName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := storage.ExpandMetaData(metaDataRaw)

	resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Queue %q (Account %s): %s", queueName, accountName, err)
	}
	if resourceGroup == nil {
		return fmt.Errorf("Unable to locate Resource Group for Storage Queue %q (Account %s)", queueName, accountName)
	}

	queueClient, err := storageClient.QueuesClient(ctx, *resourceGroup, accountName)
	if err != nil {
		return fmt.Errorf("Error building Queues Client: %s", err)
	}

	resourceID := queueClient.GetResourceID(accountName, queueName)
	if features.ShouldResourcesBeImported() {
		existing, err := queueClient.GetMetaData(ctx, accountName, queueName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Queue %q (Storage Account %q): %s", queueName, accountName, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_storage_queue", resourceID)
		}
	}

	if _, err := queueClient.Create(ctx, accountName, queueName, metaData); err != nil {
		return fmt.Errorf("Error creating Queue %q (Account %q): %+v", queueName, accountName, err)
	}

	d.SetId(resourceID)

	return resourceArmStorageQueueRead(d, meta)
}

func resourceArmStorageQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).storage
	ctx := meta.(*ArmClient).StopContext

	id, err := queues.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := storage.ExpandMetaData(metaDataRaw)

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Queue %q (Account %s): %s", id.QueueName, id.AccountName, err)
	}
	if resourceGroup == nil {
		return fmt.Errorf("Unable to locate Resource Group for Storage Queue %q (Account %s)", id.QueueName, id.AccountName)
	}

	queuesClient, err := storageClient.QueuesClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Queues Client: %s", err)
	}

	if _, err := queuesClient.SetMetaData(ctx, id.AccountName, id.QueueName, metaData); err != nil {
		return fmt.Errorf("Error setting MetaData for Queue %q (Storage Account %q): %s", id.QueueName, id.AccountName, err)
	}

	return resourceArmStorageQueueRead(d, meta)
}

func resourceArmStorageQueueRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).storage
	ctx := meta.(*ArmClient).StopContext

	id, err := queues.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Queue Container %q (Account %s): %s", id.QueueName, id.AccountName, err)
	}
	if resourceGroup == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Queue %q (Account %s) - assuming removed & removing from state", id.QueueName, id.AccountName)
		d.SetId("")
		return nil
	}

	queuesClient, err := storageClient.QueuesClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Queues Client: %s", err)
	}

	metaData, err := queuesClient.GetMetaData(ctx, id.AccountName, id.QueueName)
	if err != nil {
		if utils.ResponseWasNotFound(metaData.Response) {
			log.Printf("[INFO] Storage Queue %q no longer exists, removing from state...", id.QueueName)
			d.SetId("")
			return nil
		}

		return nil
	}

	d.Set("name", id.QueueName)
	d.Set("storage_account_name", id.AccountName)
	d.Set("resource_group_name", resourceGroup)

	if err := d.Set("metadata", storage.FlattenMetaData(metaData.MetaData)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %s", err)
	}

	return nil
}

func resourceArmStorageQueueDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).storage
	ctx := meta.(*ArmClient).StopContext

	id, err := queues.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Queue %q (Account %s): %s", id.QueueName, id.AccountName, err)
	}
	if resourceGroup == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Queue %q (Account %s) - assuming removed & removing from state", id.QueueName, id.AccountName)
		d.SetId("")
		return nil
	}

	queuesClient, err := storageClient.QueuesClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Queues Client: %s", err)
	}

	if _, err := queuesClient.Delete(ctx, id.AccountName, id.QueueName); err != nil {
		return fmt.Errorf("Error deleting Storage Queue %q: %s", id.QueueName, err)
	}

	return nil
}
