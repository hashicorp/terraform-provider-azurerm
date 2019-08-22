package azurerm

import (
	"fmt"
	"log"
	"strings"

	legacy "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/blobs"
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
				// TODO: add validation
			},

			// TODO: this can be deprecated with the new sdk?
			"resource_group_name": azure.SchemaResourceGroupName(),

			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: add validation
			},

			"storage_container_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: add validation
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

			"metadata": storage.MetaDataSchema(),
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
		options := &legacy.CopyOptions{}
		if err := blob.Copy(sourceUri, options); err != nil {
			return fmt.Errorf("Error creating storage blob on Azure: %s", err)
		}
	} else {
		switch strings.ToLower(blobType) {
		case "block":
			options := &legacy.PutBlobOptions{}
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
				options := &legacy.PutBlobOptions{}

				blob.Properties.ContentLength = size
				blob.Properties.ContentType = contentType
				if err := blob.PutPageBlob(options); err != nil {
					return fmt.Errorf("Error creating storage blob on Azure: %s", err)
				}
			}
		}
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	blob.Metadata = storage.ExpandMetaData(metaDataRaw)

	opts := &legacy.SetBlobMetadataOptions{}
	if err := blob.SetMetadata(opts); err != nil {
		return fmt.Errorf("Error setting metadata for storage blob on Azure: %s", err)
	}

	d.SetId(id)
	return resourceArmStorageBlobRead(d, meta)
}

func resourceArmStorageBlobUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := blobs.ParseResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing %q: %s", d.Id(), err)
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.AccountName, armClient)
	if err != nil {
		return err
	}

	if resourceGroup == nil {
		return fmt.Errorf("Unable to determine Resource Group for Storage Account %q", id.AccountName)
	}

	blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error getting storage account %s: %+v", id.AccountName, err)
	}
	if !accountExists {
		return fmt.Errorf("Storage account %s not found in resource group %s", id.AccountName, *resourceGroup)
	}

	container := blobClient.GetContainerReference(id.ContainerName)
	blob := container.GetBlobReference(id.BlobName)

	if d.HasChange("content_type") {
		blob.Properties.ContentType = d.Get("content_type").(string)
	}

	options := &legacy.SetBlobPropertiesOptions{}
	err = blob.SetProperties(options)
	if err != nil {
		return fmt.Errorf("Error setting properties of blob %s (container %s, storage account %s): %+v", id.BlobName, id.ContainerName, id.AccountName, err)
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		blob.Metadata = storage.ExpandMetaData(metaDataRaw)

		opts := &legacy.SetBlobMetadataOptions{}
		if err := blob.SetMetadata(opts); err != nil {
			return fmt.Errorf("Error setting metadata for storage blob on Azure: %s", err)
		}
	}

	return nil
}

func resourceArmStorageBlobRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	id, err := blobs.ParseResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing %q: %s", d.Id(), err)
	}

	resourceGroup, err := determineResourceGroupForStorageAccount(id.AccountName, armClient)
	if err != nil {
		return err
	}

	if resourceGroup == nil {
		return fmt.Errorf("Unable to determine Resource Group for Storage Account %q", id.AccountName)
	}

	blobClient, accountExists, err := armClient.storage.LegacyBlobClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing blob %q from state", id.AccountName, d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Checking for existence of storage blob %q in container %q.", id.BlobName, id.ContainerName)
	container := blobClient.GetContainerReference(id.ContainerName)
	blob := container.GetBlobReference(id.BlobName)
	exists, err := blob.Exists()
	if err != nil {
		return fmt.Errorf("error checking for existence of storage blob %q: %s", id.BlobName, err)
	}

	if !exists {
		log.Printf("[INFO] Storage blob %q no longer exists, removing from state...", id.BlobName)
		d.SetId("")
		return nil
	}

	options := &legacy.GetBlobPropertiesOptions{}
	err = blob.GetProperties(options)
	if err != nil {
		return fmt.Errorf("Error getting properties of blob %s (container %s, storage account %s): %+v", id.BlobName, id.ContainerName, id.AccountName, err)
	}

	metadataOptions := &legacy.GetBlobMetadataOptions{}
	err = blob.GetMetadata(metadataOptions)
	if err != nil {
		return fmt.Errorf("Error getting metadata of blob %s (container %s, storage account %s): %+v", id.BlobName, id.ContainerName, id.AccountName, err)
	}

	d.Set("name", id.BlobName)
	d.Set("storage_container_name", id.ContainerName)
	d.Set("storage_account_name", id.AccountName)
	d.Set("resource_group_name", resourceGroup)

	d.Set("content_type", blob.Properties.ContentType)

	d.Set("source_uri", blob.Properties.CopySource)

	blobType := strings.ToLower(strings.Replace(string(blob.Properties.BlobType), "Blob", "", 1))
	d.Set("type", blobType)

	u := blob.GetURL()
	if u == "" {
		log.Printf("[INFO] URL for %q is empty", id.BlobName)
	}
	d.Set("url", u)

	if err := d.Set("metadata", storage.FlattenMetaData(blob.Metadata)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}

	return nil
}

func resourceArmStorageBlobDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).storage
	ctx := meta.(*ArmClient).StopContext

	id, err := blobs.ParseResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing %q: %s", d.Id(), err)
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Account %q: %s", id.AccountName, err)
	}
	if resourceGroup == nil {
		return fmt.Errorf("Unable to locate Resource Group for Storage Account %q (Disk %q)!", id.AccountName, uri)
	}

	blobsClient, err := storageClient.BlobsClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Blobs Client: %s", err)
	}

	log.Printf("[INFO] Deleting Blob %q from Container %q / Storage Account %q", id.BlobName, id.ContainerName, id.AccountName)
	input := blobs.DeleteInput{
		DeleteSnapshots: true,
	}
	if _, err := blobsClient.Delete(ctx, id.AccountName, id.ContainerName, id.BlobName, input); err != nil {
		return fmt.Errorf("Error deleting Blob %q (Container %q / Account %q): %s", id.BlobName, id.ContainerName, id.AccountName, err)
	}

	return nil
}
