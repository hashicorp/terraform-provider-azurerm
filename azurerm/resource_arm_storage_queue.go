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

func resourceArmStorageQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageQueueCreate,
		Read:   resourceArmStorageQueueRead,
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
			"resource_group_name": resourceGroupNameSchema(),
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext
	environment := armClient.environment

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	queueClient, accountExists, err := armClient.getQueueServiceClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", storageAccountName)
	}

	queueReference := queueClient.GetQueueReference(name)
	id := fmt.Sprintf("https://%s.queue.%s/%s", storageAccountName, environment.StorageEndpointSuffix, name)
	if requireResourcesToBeImported {
		exists, e := queueReference.Exists()
		if e != nil {
			return fmt.Errorf("Error checking if Queue %q exists (Account %q / Resource Group %q): %s", name, storageAccountName, resourceGroupName, e)
		}

		if exists {
			return tf.ImportAsExistsError("azurerm_storage_queue", id)
		}
	}

	log.Printf("[INFO] Creating queue %q in storage account %q", name, storageAccountName)
	options := &storage.QueueServiceOptions{}
	err = queueReference.Create(options)
	if err != nil {
		return fmt.Errorf("Error creating storage queue on Azure: %s", err)
	}

	d.SetId(id)
	return resourceArmStorageQueueRead(d, meta)
}

func resourceArmStorageQueueRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := parseStorageQueueID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	if err != nil {
		return err
	}

	if resourceGroup == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Account %q (assuming removed) - removing from state", id.storageAccountName)
		d.SetId("")
		return nil
	}

	queueClient, accountExists, err := armClient.getQueueServiceClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing queue %q from state", id.storageAccountName, id.queueName)
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Checking for existence of storage queue %q.", id.queueName)
	queueReference := queueClient.GetQueueReference(id.queueName)
	exists, err := queueReference.Exists()
	if err != nil {
		return fmt.Errorf("error checking if storage queue %q exists: %s", id.queueName, err)
	}

	if !exists {
		log.Printf("[INFO] Storage queue %q no longer exists, removing from state...", id.queueName)
		d.SetId("")
		return nil
	}

	d.Set("name", id.queueName)
	d.Set("storage_account_name", id.storageAccountName)
	d.Set("resource_group_name", *resourceGroup)

	return nil
}

func resourceArmStorageQueueDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := parseStorageQueueID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	if err != nil {
		return err
	}

	if resourceGroup == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Account %q (assuming removed) - removing from state", id.storageAccountName)
		return nil
	}

	queueClient, accountExists, err := armClient.getQueueServiceClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[INFO]Storage Account %q doesn't exist so the blob won't exist", id.storageAccountName)
		return nil
	}

	log.Printf("[INFO] Deleting storage queue %q", id.queueName)
	queueReference := queueClient.GetQueueReference(id.queueName)
	options := &storage.QueueServiceOptions{}
	if err = queueReference.Delete(options); err != nil {
		return fmt.Errorf("Error deleting storage queue %q: %s", id.queueName, err)
	}

	return nil
}

type storageQueueId struct {
	storageAccountName string
	queueName          string
}

func parseStorageQueueID(input string) (*storageQueueId, error) {
	// https://myaccount.queue.core.windows.net/myqueue
	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("Error parsing %q as a URI: %+v", input, err)
	}

	segments := strings.Split(uri.Host, ".")
	if len(segments) > 0 {
		storageAccountName := segments[0]
		// remove the leading `/`
		queue := strings.TrimPrefix(uri.Path, "/")
		id := storageQueueId{
			storageAccountName: storageAccountName,
			queueName:          queue,
		}
		return &id, nil
	}

	return nil, nil
}
