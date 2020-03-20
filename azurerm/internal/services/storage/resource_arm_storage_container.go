package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/containers"
)

func resourceArmStorageContainer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceArmStorageContainerCreate,
		Read:          resourceArmStorageContainerRead,
		Delete:        resourceArmStorageContainerDelete,
		Update:        resourceArmStorageContainerUpdate,
		MigrateState:  ResourceStorageContainerMigrateState,
		SchemaVersion: 1,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				ValidateFunc: validate.StorageContainerName,
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateArmStorageAccountName,
			},

			"container_access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
				ValidateFunc: validation.StringInSlice([]string{
					"blob",
					"container",
					"private",
				}, false),
			},

			"metadata": MetaDataComputedSchema(),

			// TODO: support for ACL's, Legal Holds and Immutability Policies
			"has_immutability_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"has_legal_hold": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"resource_manager_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmStorageContainerCreate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	accessLevelRaw := d.Get("container_access_type").(string)
	azureAccessLevel := expandAzureStorageContainerAccessLevel(accessLevelRaw)
	giovanniAccessLevel := expandGiovanniStorageContainerAccessLevel(accessLevelRaw)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	azureMetaData := expandAzureMetaData(metaDataRaw)
	giovanniMetaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", accountName, containerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	azureClient := storageClient.BlobContainersClient
	giovanniClient, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Containers Client: %s", err)
	}

	id := giovanniClient.GetResourceID(accountName, containerName)

	if features.ShouldResourcesBeImported() {
		if storageClient.StorageUseAzureAD {
			err := checkContainerExistenceByAzure(azureClient, ctx, account.ResourceGroup, accountName, containerName, id)
			if err != nil {
				return err
			}
		} else {
			// When checking with giovanni and failed, give Azure client a try.
			// TODO
			// Do not check it when error is an import error.
			err := checkContainerExistenceByGiovanni(*giovanniClient, ctx, account.ResourceGroup, accountName, containerName, id)
			if err != nil {
				err := checkContainerExistenceByAzure(azureClient, ctx, account.ResourceGroup, accountName, containerName, id)
				if err != nil {
					return err
				}
			}
		}

	}

	log.Printf("[INFO] Creating Container %q in Storage Account %q", containerName, accountName)
	blobContainer := storage.BlobContainer{
		ContainerProperties: &storage.ContainerProperties{
			PublicAccess: azureAccessLevel,
			Metadata:     azureMetaData,
		},
	}
	input := containers.CreateInput{
		AccessLevel: giovanniAccessLevel,
		MetaData:    giovanniMetaData,
	}

	if storageClient.StorageUseAzureAD {
		err := createBlobContainerByAzure(ctx, azureClient, account.ResourceGroup, accountName, containerName, blobContainer)
		if err != nil {
			return err
		}
	} else {
		// When creating with giovanni and failed, give Azure client a try.
		err := createBlobContainerByGiovanni(ctx, giovanniClient, accountName, containerName, input)
		if err != nil {
			err := createBlobContainerByAzure(ctx, azureClient, account.ResourceGroup, accountName, containerName, blobContainer)
			if err != nil {
				return err
			}
		}
	}

	d.SetId(id)
	return resourceArmStorageContainerRead(d, meta)
}

func resourceArmStorageContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", id.AccountName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	accessLevelRaw := d.Get("container_access_type").(string)
	metaDataRaw := d.Get("metadata").(map[string]interface{})

	azureClient := storageClient.BlobContainersClient
	giovanniClient, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Containers Client: %s", err)
	}

	if storageClient.StorageUseAzureAD {
		log.Printf("[DEBUG] Computing the Azure Access Control for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
		accessLevel := expandAzureStorageContainerAccessLevel(accessLevelRaw)

		log.Printf("[DEBUG] Computing the Azure MetaData for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
		metaData := expandAzureMetaData(metaDataRaw)

		blobContainer := storage.BlobContainer{
			ContainerProperties: &storage.ContainerProperties{
				PublicAccess: accessLevel,
				Metadata:     metaData,
			},
		}

		err := updateBlobContainerByAzure(ctx, azureClient, account.ResourceGroup, id.AccountName, id.ContainerName, blobContainer, d)

		if err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] Computing the Giovanni Access Control for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
		accessLevel := expandGiovanniStorageContainerAccessLevel(accessLevelRaw)

		log.Printf("[DEBUG] Computing the Giovanni MetaData for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
		metaData := ExpandMetaData(metaDataRaw)

		err := updateBlobContainerByGiovanni(ctx, *giovanniClient, account.ResourceGroup, id.AccountName, id.ContainerName, metaData, accessLevel, d)
		if err != nil {
			log.Printf("[DEBUG] Computing the Azure Access Control for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
			accessLevel := expandAzureStorageContainerAccessLevel(accessLevelRaw)

			log.Printf("[DEBUG] Computing the Azure MetaData for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
			metaData := expandAzureMetaData(metaDataRaw)

			blobContainer := storage.BlobContainer{
				ContainerProperties: &storage.ContainerProperties{
					PublicAccess: accessLevel,
					Metadata:     metaData,
				},
			}

			err := updateBlobContainerByAzure(ctx, azureClient, account.ResourceGroup, id.AccountName, id.ContainerName, blobContainer, d)

			if err != nil {
				return err
			}
		}
	}

	return resourceArmStorageContainerRead(d, meta)
}

func resourceArmStorageContainerRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", id.AccountName, id.ContainerName, err)
	}
	if account == nil {
		log.Printf("[DEBUG] Unable to locate Account %q for Storage Container %q - assuming removed & removing from state", id.AccountName, id.ContainerName)
		d.SetId("")
		return nil
	}

	azureClient := storageClient.BlobContainersClient
	giovanniClient, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Containers Client: %s", err)
	}

	if storageClient.StorageUseAzureAD {
		err := readBlobContainerByAzure(ctx, azureClient, account.ResourceGroup, id.AccountName, id.ContainerName, d)
		if err != nil {
			return err
		}
	} else {
		// Give Azure a try when giovanni fails
		err := readBlobContainerByGiovanni(ctx, *giovanniClient, account.ResourceGroup, id.AccountName, id.ContainerName, d)
		if err != nil {
			err := readBlobContainerByAzure(ctx, azureClient, account.ResourceGroup, id.AccountName, id.ContainerName, d)
			if err != nil {
				return err
			}
		}
	}

	d.Set("name", id.ContainerName)
	d.Set("storage_account_name", id.AccountName)

<<<<<<< HEAD
	accessLevel := flattenAzureStorageContainerAccessLevel(props.PublicAccess)

	d.Set("container_access_type", accessLevel)

	if err := d.Set("metadata", flattenAzureMetaData(props.Metadata)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}

	d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
	d.Set("has_legal_hold", props.HasLegalHold)

	resourceManagerId := client.GetResourceManagerResourceID(storageClient.SubscriptionId, account.ResourceGroup, id.AccountName, id.ContainerName)
	d.Set("resource_manager_id", resourceManagerId)

=======
>>>>>>> Read props by Giovanni and/or Azure
	return nil
}

func resourceArmStorageContainerDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", id.AccountName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	azureClient := storageClient.BlobContainersClient

	if _, err := azureClient.Delete(ctx, account.ResourceGroup, id.AccountName, id.ContainerName); err != nil {
		return fmt.Errorf("Error deleting Container %q (Storage Account %q / Resource Group %q): %s", id.ContainerName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

func getAzureBlobContainerProperties(accessLevelRaw string, metaDataRaw map[string]interface{}) (*storage.BlobContainer, error) {
	// For backward compatibility, raw access value has to be converted.
	// expandAzureStorageContainerAccessLevel will an empty string if it cannot find a value
	// that maps to a storage.PublicAccess value.
	// Therefore, if parsed value is an empty string, we are facing an error.
	// It does not seem to be a good way, but it is the cost to use switch.
	accessLevel := expandAzureStorageContainerAccessLevel(accessLevelRaw)
	if string(accessLevel) == "" {
		return nil, fmt.Errorf("Error parse %q to a Azure blob container access level")
	}

	metaData := expandAzureMetaData(metaDataRaw)

	return &storage.BlobContainer{
		ContainerProperties: &storage.ContainerProperties{
			PublicAccess: accessLevel,
			Metadata:     metaData,
		},
	}, nil
}

func getGiovanniBlobContainerProperties(accessLevelRaw string, metaDataRaw map[string]interface{}) containers.CreateInput {
	accessLevel := expandGiovanniStorageContainerAccessLevel(accessLevelRaw)

	metaData := ExpandMetaData(metaDataRaw)

	return containers.CreateInput{
		AccessLevel: accessLevel,
		MetaData:    metaData,
	}
}

func expandAzureStorageContainerAccessLevel(input string) storage.PublicAccess {
	switch input {
	case "private":
		return storage.PublicAccessNone
	case "container":
		return storage.PublicAccessContainer
	case "blob":
		return storage.PublicAccessBlob
	default:
		return storage.PublicAccess("")
	}
}

func flattenAzureStorageContainerAccessLevel(input storage.PublicAccess) string {
	switch input {
	case storage.PublicAccessNone:
		return "private"
	case storage.PublicAccessContainer:
		return "container"
	case storage.PublicAccessBlob:
		return "blob"
	default:
		return string(input)
	}
}

func getAzureResourceID(baseUri, accountName, containerName string) string {
	// For backforward compatible, generate resource ID in the same way as giovanni's.
	domain := parsers.GetBlobEndpoint(baseUri, accountName)
	return fmt.Sprintf("%s/%s", domain, containerName)
}

func expandAzureMetaData(input map[string]interface{}) map[string]*string {
	output := make(map[string]*string)

	for k, v := range input {
		temp := v.(string)
		output[k] = &temp
	}

	return output
}

func flattenAzureMetaData(input map[string]*string) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		output[k] = v
	}

	return output
}

func expandGiovanniStorageContainerAccessLevel(input string) containers.AccessLevel {
	// for historical reasons, "private" above is an empty string in the API
	// so the enum doesn't 1:1 match. You could argue the SDK should handle this
	// but this is suitable for now
	if input == "private" {
		return containers.Private
	}

	return containers.AccessLevel(input)
}

func flattenGiovanniStorageContainerAccessLevel(input containers.AccessLevel) string {
	// for historical reasons, "private" above is an empty string in the API
	if input == containers.Private {
		return "private"
	}

	return string(input)
}

func checkContainerExistenceByAzure(azClient storage.BlobContainersClient, ctx context.Context, resourceGroup, accountName, containerName, id string) error {
	existing, err := azClient.Get(ctx, resourceGroup, accountName, containerName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for existence of existing Container %q (Account %q / Resource Group %q) with Azure BlobContainersClient: %+v", containerName, accountName, resourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_storage_container", id)
	}

	return nil
}

func checkContainerExistenceByGiovanni(gvnClient containers.Client, ctx context.Context, resourceGroup, accountName, containerName, id string) error {
	existing, err := gvnClient.GetProperties(ctx, accountName, containerName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for existence of existing Container %q (Account %q / Resource Group %q) with Giovanni ContainersClient: %+v", containerName, accountName, resourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_storage_container", id)
	}
	return nil
}

func createBlobContainerByAzure(ctx context.Context, azClient storage.BlobContainersClient, resourceGroup, accountName, containerName string, blobContainer storage.BlobContainer) error {
	_, err := azClient.Create(ctx, resourceGroup, accountName, containerName, blobContainer)

	return err
}

func createBlobContainerByGiovanni(ctx context.Context, gvnClient *containers.Client, accountName, containerName string, input containers.CreateInput) error {
	_, err := gvnClient.Create(ctx, accountName, containerName, input)

	return err
}

func readBlobContainerByAzure(ctx context.Context, azClient storage.BlobContainersClient, resourceGroup, accountName, containerName string, d *schema.ResourceData) error {
	props, err := azClient.Get(ctx, resourceGroup, accountName, containerName)
	if err != nil {
		if utils.ResponseWasNotFound(props.Response) {
			log.Printf("[DEBUG] Container %q was not found in Account %q / Resource Group %q - assuming removed & removing from state", containerName, accountName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Container %q (Account %q / Resource Group %q): %s", containerName, accountName, resourceGroup, err)
	}

	d.Set("container_access_type", flattenAzureStorageContainerAccessLevel(props.PublicAccess))

	if err := d.Set("metadata", flattenAzureMetaData(props.Metadata)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}
	d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
	d.Set("has_legal_hold", props.HasLegalHold)
	return nil
}

func readBlobContainerByGiovanni(ctx context.Context, gvnClient containers.Client, resourceGroup, accountName, containerName string, d *schema.ResourceData) error {
	props, err := gvnClient.GetProperties(ctx, accountName, containerName)
	if err != nil {
		if utils.ResponseWasNotFound(props.Response) {
			log.Printf("[DEBUG] Container %q was not found in Account %q / Resource Group %q - assuming removed & removing from state", containerName, accountName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Container %q (Account %q / Resource Group %q): %s", containerName, accountName, resourceGroup, err)
	}

	d.Set("container_access_type", flattenGiovanniStorageContainerAccessLevel(props.AccessLevel))

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}
	d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
	d.Set("has_legal_hold", props.HasLegalHold)
	return nil
}

func updateBlobContainerByAzure(ctx context.Context, azClient storage.BlobContainersClient, resourceGroup, accountName, containerName string, blobContainer storage.BlobContainer, d *schema.ResourceData) error {
	return nil
}

func updateBlobContainerByGiovanni(ctx context.Context, gvnClient containers.Client, resourceGroup, accountName, containerName string, metaData map[string]string, accessLevel containers.AccessLevel, d *schema.ResourceData) error {
	return nil
}
