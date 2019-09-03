package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
				// TODO: @tombuildsstuff - a note this only works for Page blobs
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      8,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"metadata": storage.MetaDataSchema(),

			// Deprecated fields
			"attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				Deprecated:   "Retries are now handled by the Azure SDK as such this field is no longer necessary and will be removed in v2.0 of the Azure Provider",
				ValidateFunc: validation.IntAtLeast(1),
			},

			"resource_group_name": azure.SchemaResourceGroupNameDeprecated(),
		},
	}
}

func resourceArmStorageBlobCreate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*ArmClient).storage
	ctx := meta.(*ArmClient).StopContext

	accountName := d.Get("storage_account_name").(string)
	containerName := d.Get("storage_container_name").(string)
	name := d.Get("name").(string)

	resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Account %q: %s", accountName, err)
	}
	if resourceGroup == nil {
		return fmt.Errorf("Unable to locate Resource Group for Blob %q (Container %q / Account %q)", name, containerName, accountName)
	}

	blobsClient, err := storageClient.BlobsClient(ctx, *resourceGroup, accountName)
	if err != nil {
		return fmt.Errorf("Error building Blobs Client: %s", err)
	}

	id := blobsClient.GetResourceID(accountName, containerName, name)
	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		input := blobs.GetPropertiesInput{}
		props, err := blobsClient.GetProperties(ctx, accountName, containerName, name, input)
		if err != nil {
			if !utils.ResponseWasNotFound(props.Response) {
				return fmt.Errorf("Error checking if Blob %q exists (Container %q / Account %q / Resource Group %q): %s", name, containerName, accountName, *resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(props.Response) {
			return tf.ImportAsExistsError("azurerm_storage_blob", id)
		}
	}

	log.Printf("[DEBUG] Creating Blob %q in Container %q within Storage Account %q..", name, containerName, accountName)
	metaDataRaw := d.Get("metadata").(map[string]interface{})
	blobInput := storage.BlobUpload{
		AccountName:   accountName,
		ContainerName: containerName,
		BlobName:      name,
		Client:        blobsClient,

		BlobType:    d.Get("type").(string),
		ContentType: d.Get("content_type").(string),
		MetaData:    storage.ExpandMetaData(metaDataRaw),
		Parallelism: d.Get("parallelism").(int),
		Size:        d.Get("size").(int),
		Source:      d.Get("source").(string),
		SourceUri:   d.Get("source_uri").(string),
	}
	if err := blobInput.Create(ctx); err != nil {
		return fmt.Errorf("Error creating Blob %q (Container %q / Account %q): %s", name, containerName, accountName, err)
	}
	log.Printf("[DEBUG] Created Blob %q in Container %q within Storage Account %q.", name, containerName, accountName)

	d.SetId(id)

	return resourceArmStorageBlobUpdate(d, meta)
}

func resourceArmStorageBlobUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Unable to locate Resource Group for Blob %q (Container %q / Account %q)", id.BlobName, id.ContainerName, id.AccountName)
	}

	blobsClient, err := storageClient.BlobsClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Blobs Client: %s", err)
	}

	// TODO: changing the access tier

	if d.HasChange("content_type") {
		log.Printf("[DEBUG] Updating Properties for Blob %q (Container %q / Account %q)...", id.BlobName, id.ContainerName, id.AccountName)
		input := blobs.SetPropertiesInput{
			ContentType: utils.String(d.Get("content_type").(string)),
		}
		if _, err := blobsClient.SetProperties(ctx, id.AccountName, id.ContainerName, id.BlobName, input); err != nil {
			return fmt.Errorf("Error updating Properties for Blob %q (Container %q / Account %q): %s", id.BlobName, id.ContainerName, id.AccountName, err)
		}
		log.Printf("[DEBUG] Updated Properties for Blob %q (Container %q / Account %q).", id.BlobName, id.ContainerName, id.AccountName)
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating MetaData for Blob %q (Container %q / Account %q)...", id.BlobName, id.ContainerName, id.AccountName)
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		input := blobs.SetMetaDataInput{
			MetaData: storage.ExpandMetaData(metaDataRaw),
		}
		if _, err := blobsClient.SetMetaData(ctx, id.AccountName, id.ContainerName, id.BlobName, input); err != nil {
			return fmt.Errorf("Error updating MetaData for Blob %q (Container %q / Account %q): %s", id.BlobName, id.ContainerName, id.AccountName, err)
		}
		log.Printf("[DEBUG] Updated MetaData for Blob %q (Container %q / Account %q).", id.BlobName, id.ContainerName, id.AccountName)
	}

	return nil
}

func resourceArmStorageBlobRead(d *schema.ResourceData, meta interface{}) error {
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
		log.Printf("[DEBUG] Unable to locate Resource Group for Blob %q (Container %q / Account %q) - assuming removed & removing from state!", id.BlobName, id.ContainerName, id.AccountName)
		d.SetId("")
		return nil
	}

	blobsClient, err := storageClient.BlobsClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Blobs Client: %s", err)
	}

	log.Printf("[INFO] Retrieving Storage Blob %q (Container %q / Account %q).", id.BlobName, id.ContainerName, id.AccountName)
	input := blobs.GetPropertiesInput{}
	props, err := blobsClient.GetProperties(ctx, id.AccountName, id.ContainerName, id.BlobName, input)
	if err != nil {
		if utils.ResponseWasNotFound(props.Response) {
			log.Printf("[INFO] Blob %q was not found in Container %q / Account %q - assuming removed & removing from state...", id.BlobName, id.ContainerName, id.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving properties for Blob %q (Container %q / Account %q): %s", id.BlobName, id.ContainerName, id.AccountName, err)
	}

	d.Set("name", id.BlobName)
	d.Set("storage_container_name", id.ContainerName)
	d.Set("storage_account_name", id.AccountName)
	d.Set("resource_group_name", resourceGroup)

	d.Set("content_type", props.ContentType)

	// The CopySource is only returned if the blob hasn't been modified (e.g. metadata configured etc)
	// as such, we need to conditionally set this to ensure it's trackable if possible
	if props.CopySource != "" {
		d.Set("source_uri", props.CopySource)
	}

	blobType := strings.ToLower(strings.Replace(string(props.BlobType), "Blob", "", 1))
	d.Set("type", blobType)

	d.Set("url", d.Id())

	if err := d.Set("metadata", storage.FlattenMetaData(props.MetaData)); err != nil {
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
		return fmt.Errorf("Unable to locate Resource Group for Storage Account %q!", id.AccountName)
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
