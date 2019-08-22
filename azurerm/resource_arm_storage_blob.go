package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmStorageBlob() *schema.Resource {
	return &schema.Resource{
		Create:        resourceArmStorageBlobCreate,
		Read:          resourceArmStorageBlobRead,
		Update:        resourceArmStorageBlobUpdate,
		Delete:        resourceArmStorageBlobDelete,
		MigrateState:  resourceStorageBlobMigrateState,
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"storage_container_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"block", "page"}, true),
			},

			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      0,
				ValidateFunc: validate.IntDivisibleBy(512),
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "application/octet-stream",
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_uri"},
			},

			"source_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source"},
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"parallelism": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      8,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
			},
		},
	}
}

func resourceArmStorageBlobCreate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext
	env := armClient.environment

	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", storageAccountName)
	}

	name := d.Get("name").(string)
	blobType := d.Get("type").(string)
	containerName := d.Get("storage_container_name").(string)
	sourceUri := d.Get("source_uri").(string)
	contentType := d.Get("content_type").(string)

	log.Printf("[INFO] Creating blob %q in container %q within storage account %q", name, containerName, storageAccountName)
	container := blobClient.GetContainerReference(containerName)
	blob := container.GetBlobReference(name)

	// gives us https://example.blob.core.windows.net/container/file.vhd
	id := fmt.Sprintf("https://%s.blob.%s/%s/%s", storageAccountName, env.StorageEndpointSuffix, containerName, name)
	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		exists, err := blob.Exists()
		if err != nil {
			return fmt.Errorf("Error checking if Blob %q exists (Container %q / Account %q / Resource Group %q): %s", name, containerName, storageAccountName, resourceGroupName, err)
		}

		if exists {
			return tf.ImportAsExistsError("azurerm_storage_blob", id)
		}
	}

	if sourceUri != "" {
		options := &storage.CopyOptions{}
		if err := blob.Copy(sourceUri, options); err != nil {
			return fmt.Errorf("Error creating storage blob on Azure: %s", err)
		}
	} else {
		switch strings.ToLower(blobType) {
		case "block":
			options := &storage.PutBlobOptions{}
			if err := blob.CreateBlockBlob(options); err != nil {
				return fmt.Errorf("Error creating storage blob on Azure: %s", err)
			}

			source := d.Get("source").(string)
			if source != "" {
				parallelism := d.Get("parallelism").(int)
				attempts := d.Get("attempts").(int)

				if err := resourceArmStorageBlobBlockUploadFromSource(containerName, name, source, contentType, blobClient, parallelism, attempts); err != nil {
					return fmt.Errorf("Error creating storage blob on Azure: %s", err)
				}
			}
		case "page":
			source := d.Get("source").(string)
			if source != "" {
				parallelism := d.Get("parallelism").(int)
				attempts := d.Get("attempts").(int)

				if err := resourceArmStorageBlobPageUploadFromSource(containerName, name, source, contentType, blobClient, parallelism, attempts); err != nil {
					return fmt.Errorf("Error creating storage blob on Azure: %s", err)
				}
			} else {
				size := int64(d.Get("size").(int))
				options := &storage.PutBlobOptions{}

				blob.Properties.ContentLength = size
				blob.Properties.ContentType = contentType
				if err := blob.PutPageBlob(options); err != nil {
					return fmt.Errorf("Error creating storage blob on Azure: %s", err)
				}
			}
		}
	}

	blob.Metadata = expandStorageAccountBlobMetadata(d)

	opts := &storage.SetBlobMetadataOptions{}
	if err := blob.SetMetadata(opts); err != nil {
		return fmt.Errorf("Error setting metadata for storage blob on Azure: %s", err)
	}

	d.SetId(id)
	return resourceArmStorageBlobRead(d, meta)
}

func resourceArmStorageBlobUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := parseStorageBlobID(d.Id(), armClient.environment)
	if err != nil {
		return err
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	if err != nil {
		return err
	}

	if resourceGroup == nil {
		return fmt.Errorf("Unable to determine Resource Group for Storage Account %q", id.storageAccountName)
	}

	blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return fmt.Errorf("Error getting storage account %s: %+v", id.storageAccountName, err)
	}
	if !accountExists {
		return fmt.Errorf("Storage account %s not found in resource group %s", id.storageAccountName, *resourceGroup)
	}

	container := blobClient.GetContainerReference(id.containerName)
	blob := container.GetBlobReference(id.blobName)

	if d.HasChange("content_type") {
		blob.Properties.ContentType = d.Get("content_type").(string)
	}

	options := &storage.SetBlobPropertiesOptions{}
	err = blob.SetProperties(options)
	if err != nil {
		return fmt.Errorf("Error setting properties of blob %s (container %s, storage account %s): %+v", id.blobName, id.containerName, id.storageAccountName, err)
	}

	if d.HasChange("metadata") {
		blob.Metadata = expandStorageAccountBlobMetadata(d)

		opts := &storage.SetBlobMetadataOptions{}
		if err := blob.SetMetadata(opts); err != nil {
			return fmt.Errorf("Error setting metadata for storage blob on Azure: %s", err)
		}
	}

	return nil
}

func resourceArmStorageBlobRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := parseStorageBlobID(d.Id(), armClient.environment)
	if err != nil {
		return err
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	if err != nil {
		return err
	}

	if resourceGroup == nil {
		return fmt.Errorf("Unable to determine Resource Group for Storage Account %q", id.storageAccountName)
	}

	blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing blob %q from state", id.storageAccountName, d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Checking for existence of storage blob %q in container %q.", id.blobName, id.containerName)
	container := blobClient.GetContainerReference(id.containerName)
	blob := container.GetBlobReference(id.blobName)
	exists, err := blob.Exists()
	if err != nil {
		return fmt.Errorf("error checking for existence of storage blob %q: %s", id.blobName, err)
	}

	if !exists {
		log.Printf("[INFO] Storage blob %q no longer exists, removing from state...", id.blobName)
		d.SetId("")
		return nil
	}

	options := &storage.GetBlobPropertiesOptions{}
	err = blob.GetProperties(options)
	if err != nil {
		return fmt.Errorf("Error getting properties of blob %s (container %s, storage account %s): %+v", id.blobName, id.containerName, id.storageAccountName, err)
	}

	metadataOptions := &storage.GetBlobMetadataOptions{}
	err = blob.GetMetadata(metadataOptions)
	if err != nil {
		return fmt.Errorf("Error getting metadata of blob %s (container %s, storage account %s): %+v", id.blobName, id.containerName, id.storageAccountName, err)
	}

	d.Set("name", id.blobName)
	d.Set("storage_container_name", id.containerName)
	d.Set("storage_account_name", id.storageAccountName)
	d.Set("resource_group_name", resourceGroup)

	d.Set("content_type", blob.Properties.ContentType)

	d.Set("source_uri", blob.Properties.CopySource)

	blobType := strings.ToLower(strings.Replace(string(blob.Properties.BlobType), "Blob", "", 1))
	d.Set("type", blobType)

	u := blob.GetURL()
	if u == "" {
		log.Printf("[INFO] URL for %q is empty", id.blobName)
	}
	d.Set("url", u)
	d.Set("metadata", flattenStorageAccountBlobMetadata(blob.Metadata))

	return nil
}

func resourceArmStorageBlobDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := parseStorageBlobID(d.Id(), armClient.environment)
	if err != nil {
		return err
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.storageAccountName, armClient)
	if err != nil {
		return fmt.Errorf("Unable to determine Resource Group for Storage Account %q: %+v", id.storageAccountName, err)
	}
	if resourceGroup == nil {
		log.Printf("[INFO] Resource Group doesn't exist so the blob won't exist")
		return nil
	}

	blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, *resourceGroup, id.storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[INFO] Storage Account %q doesn't exist so the blob won't exist", id.storageAccountName)
		return nil
	}

	log.Printf("[INFO] Deleting storage blob %q", id.blobName)
	options := &storage.DeleteBlobOptions{}
	container := blobClient.GetContainerReference(id.containerName)
	blob := container.GetBlobReference(id.blobName)
	_, err = blob.DeleteIfExists(options)
	if err != nil {
		return fmt.Errorf("Error deleting storage blob %q: %s", id.blobName, err)
	}

	return nil
}

func expandStorageAccountBlobMetadata(d *schema.ResourceData) storage.BlobMetadata {
	blobMetadata := make(map[string]string)

	blobMetadataRaw := d.Get("metadata").(map[string]interface{})
	for key, value := range blobMetadataRaw {
		blobMetadata[key] = value.(string)
	}
	return blobMetadata
}

func flattenStorageAccountBlobMetadata(in storage.BlobMetadata) map[string]interface{} {
	blobMetadata := make(map[string]interface{})

	for key, value := range in {
		blobMetadata[key] = value
	}

	return blobMetadata
}
