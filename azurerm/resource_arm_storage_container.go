package azurerm

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/storage"
	azauto "github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func resourceArmStorageContainer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceArmStorageContainerCreateUpdate,
		Read:          resourceArmStorageContainerRead,
		Delete:        resourceArmStorageContainerDelete,
		Update:        resourceArmStorageContainerCreateUpdate,
		MigrateState:  resourceStorageContainerMigrateState,
		SchemaVersion: 1,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageContainerName,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "private",
				ValidateFunc: validateArmStorageContainerAccessType,
			},

			"properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

//Following the naming convention as laid out in the docs
func validateArmStorageContainerName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^\$root$|^[0-9a-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q: %q",
			k, value))
	}
	if len(value) < 3 || len(value) > 63 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 3 and 63 characters: %q", k, value))
	}
	if regexp.MustCompile(`^-`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot begin with a hyphen: %q", k, value))
	}
	return warnings, errors
}

func validateArmStorageContainerAccessType(v interface{}, _ string) (warnings []string, errors []error) {
	value := strings.ToLower(v.(string))
	validTypes := map[string]struct{}{
		"private":   {},
		"blob":      {},
		"container": {},
	}

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf("Storage container access type %q is invalid, must be %q, %q or %q", value, "private", "blob", "page"))
	}
	return warnings, errors
}

func resourceArmStorageContainerCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", storageAccountName)
	}

	var accessType storage.ContainerAccessType
	if d.Get("container_access_type").(string) == "private" {
		accessType = storage.ContainerAccessType("")
	} else {
		accessType = storage.ContainerAccessType(d.Get("container_access_type").(string))
	}

	reference := blobClient.GetContainerReference(name)
	id := fmt.Sprintf("https://%s.blob.%s/%s", storageAccountName, armClient.environment.StorageEndpointSuffix, name)
	if requireResourcesToBeImported && d.IsNewResource() {
		exists, e := reference.Exists()
		if e != nil {
			return fmt.Errorf("Error checking if Storage Container %q exists (Account %q / Resource Group %q): %s", name, storageAccountName, resourceGroupName, e)
		}

		if exists {
			return tf.ImportAsExistsError("azurerm_storage_container", id)
		}
	}

	log.Printf("[INFO] Creating container %q in storage account %q.", name, storageAccountName)
	err = resource.Retry(120*time.Second, checkContainerIsCreated(reference))
	if err != nil {
		return fmt.Errorf("Error creating container %q in storage account %q: %s", name, storageAccountName, err)
	}

	permissions := storage.ContainerPermissions{
		AccessType: accessType,
	}
	permissionOptions := &storage.SetContainerPermissionOptions{}
	err = reference.SetPermissions(permissions, permissionOptions)
	if err != nil {
		return fmt.Errorf("Error setting permissions for container %s in storage account %s: %+v", name, storageAccountName, err)
	}

	d.SetId(id)
	return resourceArmStorageContainerRead(d, meta)
}

// resourceAzureStorageContainerRead does all the necessary API calls to
// read the status of the storage container off Azure.
func resourceArmStorageContainerRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := parseStorageContainerID(d.Id(), armClient.environment)
	if err != nil {
		return err
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	if err != nil {
		return err
	}
	if resourceGroup == nil {
		log.Printf("Cannot locate Resource Group for Storage Account %q (presuming it's gone) - removing from state", id.storageAccountName)
		d.SetId("")
		return nil
	}

	blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing container %q from state", id.storageAccountName, d.Id())
		d.SetId("")
		return nil
	}

	var container *storage.Container
	listParams := storage.ListContainersParameters{
		Prefix:  id.containerName,
		Timeout: 90,
	}

	for {
		resp, err := blobClient.ListContainers(listParams)
		if err != nil {
			return fmt.Errorf("Failed to retrieve storage resp in account %q: %s", id.containerName, err)
		}

		for _, c := range resp.Containers {
			if c.Name == id.containerName {
				container = &c
				break
			}
		}

		if resp.NextMarker == "" {
			break
		}

		listParams.Marker = resp.NextMarker
	}

	if container == nil {
		log.Printf("[INFO] Storage container %q does not exist in account %q, removing from state...", id.containerName, id.storageAccountName)
		d.SetId("")
		return nil
	}

	d.Set("name", id.containerName)
	d.Set("storage_account_name", id.storageAccountName)
	d.Set("resource_group_name", resourceGroup)

	// for historical reasons, "private" above is an empty string in the API
	if container.Properties.PublicAccess == storage.ContainerAccessTypePrivate {
		d.Set("container_access_type", "private")
	} else {
		d.Set("container_access_type", string(container.Properties.PublicAccess))
	}

	output := make(map[string]interface{})

	output["last_modified"] = container.Properties.LastModified
	output["lease_status"] = container.Properties.LeaseStatus
	output["lease_state"] = container.Properties.LeaseState
	output["lease_duration"] = container.Properties.LeaseDuration

	if err := d.Set("properties", output); err != nil {
		return fmt.Errorf("Error setting `properties`: %+v", err)
	}

	return nil
}

// resourceAzureStorageContainerDelete does all the necessary API calls to
// delete a storage container off Azure.
func resourceArmStorageContainerDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := parseStorageContainerID(d.Id(), armClient.environment)
	if err != nil {
		return err
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	if err != nil {
		return err
	}
	if resourceGroup == nil {
		log.Printf("Cannot locate Resource Group for Storage Account %q (presuming it's gone) - removing from state", id.storageAccountName)
		return nil
	}

	blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[INFO] Storage Account %q doesn't exist so the container won't exist", id.storageAccountName)
		return nil
	}

	log.Printf("[INFO] Deleting storage container %q in account %q", id.containerName, id.storageAccountName)
	reference := blobClient.GetContainerReference(id.containerName)
	deleteOptions := &storage.DeleteContainerOptions{}
	if _, err := reference.DeleteIfExists(deleteOptions); err != nil {
		return fmt.Errorf("Error deleting storage container %q from storage account %q: %s", id.containerName, id.storageAccountName, err)
	}

	return nil
}

func checkContainerIsCreated(reference *storage.Container) func() *resource.RetryError {
	return func() *resource.RetryError {
		createOptions := &storage.CreateContainerOptions{}

		if _, err := reference.CreateIfNotExists(createOptions); err != nil {
			return resource.RetryableError(err)
		}

		return nil
	}
}

type storageContainerId struct {
	storageAccountName string
	containerName      string
}

func parseStorageContainerID(input string, environment azauto.Environment) (*storageContainerId, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("Error parsing %q as URI: %+v", input, err)
	}

	// remove the leading `/`
	segments := strings.Split(strings.TrimPrefix(uri.Path, "/"), "/")
	if len(segments) < 1 {
		return nil, fmt.Errorf("Expected number of segments in the path to be < 1 but got %d", len(segments))
	}

	storageAccountName := strings.Replace(uri.Host, fmt.Sprintf(".blob.%s", environment.StorageEndpointSuffix), "", 1)
	containerName := segments[0]

	id := storageContainerId{
		storageAccountName: storageAccountName,
		containerName:      containerName,
	}
	return &id, nil
}
