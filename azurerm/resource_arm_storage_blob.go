package azurerm

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/Azure/azure-sdk-for-go/storage"
	azauto "github.com/Azure/go-autorest/autorest/azure"
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

	blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
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
	if requireResourcesToBeImported && d.IsNewResource() {
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

type resourceArmStorageBlobPage struct {
	offset  int64
	section *io.SectionReader
}

func resourceArmStorageBlobPageUploadFromSource(container, name, source, contentType string, client *storage.BlobStorageClient, parallelism, attempts int) error {
	workerCount := parallelism * runtime.NumCPU()

	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("Error opening source file for upload %q: %s", source, err)
	}
	defer utils.IoCloseAndLogError(file, fmt.Sprintf("Error closing Storage Blob `%s` file `%s` after upload", name, source))

	blobSize, pageList, err := resourceArmStorageBlobPageSplit(file)
	if err != nil {
		return fmt.Errorf("Error splitting source file %q into pages: %s", source, err)
	}

	options := &storage.PutBlobOptions{}
	containerRef := client.GetContainerReference(container)
	blob := containerRef.GetBlobReference(name)
	blob.Properties.ContentLength = blobSize
	blob.Properties.ContentType = contentType
	err = blob.PutPageBlob(options)
	if err != nil {
		return fmt.Errorf("Error creating storage blob on Azure: %s", err)
	}

	pages := make(chan resourceArmStorageBlobPage, len(pageList))
	errors := make(chan error, len(pageList))
	wg := &sync.WaitGroup{}
	wg.Add(len(pageList))

	total := int64(0)
	for _, page := range pageList {
		total += page.section.Size()
		pages <- page
	}
	close(pages)

	for i := 0; i < workerCount; i++ {
		go resourceArmStorageBlobPageUploadWorker(resourceArmStorageBlobPageUploadContext{
			container: container,
			name:      name,
			source:    source,
			blobSize:  blobSize,
			client:    client,
			pages:     pages,
			errors:    errors,
			wg:        wg,
			attempts:  attempts,
		})
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("Error while uploading source file %q: %s", source, <-errors)
	}

	return nil
}

func resourceArmStorageBlobPageSplit(file *os.File) (int64, []resourceArmStorageBlobPage, error) {
	const (
		minPageSize int64 = 4 * 1024
		maxPageSize int64 = 4 * 1024 * 1024
	)

	info, err := file.Stat()
	if err != nil {
		return int64(0), nil, fmt.Errorf("Could not stat file %q: %s", file.Name(), err)
	}

	blobSize := info.Size()
	if info.Size()%minPageSize != 0 {
		blobSize = info.Size() + (minPageSize - (info.Size() % minPageSize))
	}

	emptyPage := make([]byte, minPageSize)

	type byteRange struct {
		offset int64
		length int64
	}

	var nonEmptyRanges []byteRange
	var currentRange byteRange
	for i := int64(0); i < blobSize; i += minPageSize {
		pageBuf := make([]byte, minPageSize)
		_, err = file.ReadAt(pageBuf, i)
		if err != nil && err != io.EOF {
			return int64(0), nil, fmt.Errorf("Could not read chunk at %d: %s", i, err)
		}

		if bytes.Equal(pageBuf, emptyPage) {
			if currentRange.length != 0 {
				nonEmptyRanges = append(nonEmptyRanges, currentRange)
			}
			currentRange = byteRange{
				offset: i + minPageSize,
			}
		} else {
			currentRange.length += minPageSize
			if currentRange.length == maxPageSize || (currentRange.offset+currentRange.length == blobSize) {
				nonEmptyRanges = append(nonEmptyRanges, currentRange)
				currentRange = byteRange{
					offset: i + minPageSize,
				}
			}
		}
	}

	var pages []resourceArmStorageBlobPage
	for _, nonEmptyRange := range nonEmptyRanges {
		pages = append(pages, resourceArmStorageBlobPage{
			offset:  nonEmptyRange.offset,
			section: io.NewSectionReader(file, nonEmptyRange.offset, nonEmptyRange.length),
		})
	}

	return info.Size(), pages, nil
}

type resourceArmStorageBlobPageUploadContext struct {
	container string
	name      string
	source    string
	blobSize  int64
	client    *storage.BlobStorageClient
	pages     chan resourceArmStorageBlobPage
	errors    chan error
	wg        *sync.WaitGroup
	attempts  int
}

func resourceArmStorageBlobPageUploadWorker(ctx resourceArmStorageBlobPageUploadContext) {
	for page := range ctx.pages {
		start := page.offset
		end := page.offset + page.section.Size() - 1
		if end > ctx.blobSize-1 {
			end = ctx.blobSize - 1
		}
		size := end - start + 1

		chunk := make([]byte, size)
		_, err := page.section.Read(chunk)
		if err != nil && err != io.EOF {
			ctx.errors <- fmt.Errorf("Error reading source file %q at offset %d: %s", ctx.source, page.offset, err)
			ctx.wg.Done()
			continue
		}

		for x := 0; x < ctx.attempts; x++ {
			container := ctx.client.GetContainerReference(ctx.container)
			blob := container.GetBlobReference(ctx.name)
			blobRange := storage.BlobRange{
				Start: uint64(start),
				End:   uint64(end),
			}
			options := &storage.PutPageOptions{}
			reader := bytes.NewReader(chunk)
			err = blob.WriteRange(blobRange, reader, options)
			if err == nil {
				break
			}
		}
		if err != nil {
			ctx.errors <- fmt.Errorf("Error writing page at offset %d for file %q: %s", page.offset, ctx.source, err)
			ctx.wg.Done()
			continue
		}

		ctx.wg.Done()
	}
}

type resourceArmStorageBlobBlock struct {
	section *io.SectionReader
	id      string
}

func resourceArmStorageBlobBlockUploadFromSource(container, name, source, contentType string, client *storage.BlobStorageClient, parallelism, attempts int) error {
	workerCount := parallelism * runtime.NumCPU()

	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("Error opening source file for upload %q: %s", source, err)
	}
	defer utils.IoCloseAndLogError(file, fmt.Sprintf("Error closing Storage Blob `%s` file `%s` after upload", name, source))

	blockList, parts, err := resourceArmStorageBlobBlockSplit(file)
	if err != nil {
		return fmt.Errorf("Error reading and splitting source file for upload %q: %s", source, err)
	}

	wg := &sync.WaitGroup{}
	blocks := make(chan resourceArmStorageBlobBlock, len(parts))
	errors := make(chan error, len(parts))

	wg.Add(len(parts))
	for _, p := range parts {
		blocks <- p
	}
	close(blocks)

	for i := 0; i < workerCount; i++ {
		go resourceArmStorageBlobBlockUploadWorker(resourceArmStorageBlobBlockUploadContext{
			client:    client,
			source:    source,
			container: container,
			name:      name,
			blocks:    blocks,
			errors:    errors,
			wg:        wg,
			attempts:  attempts,
		})
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("Error while uploading source file %q: %s", source, <-errors)
	}

	containerReference := client.GetContainerReference(container)
	blobReference := containerReference.GetBlobReference(name)
	blobReference.Properties.ContentType = contentType
	options := &storage.PutBlockListOptions{}
	err = blobReference.PutBlockList(blockList, options)
	if err != nil {
		return fmt.Errorf("Error updating block list for source file %q: %s", source, err)
	}

	return nil
}

func resourceArmStorageBlobBlockSplit(file *os.File) ([]storage.Block, []resourceArmStorageBlobBlock, error) {
	const (
		idSize          = 64
		blockSize int64 = 4 * 1024 * 1024
	)
	var parts []resourceArmStorageBlobBlock
	var blockList []storage.Block

	info, err := file.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("Error stating source file %q: %s", file.Name(), err)
	}

	for i := int64(0); i < info.Size(); i = i + blockSize {
		entropy := make([]byte, idSize)
		_, err = rand.Read(entropy)
		if err != nil {
			return nil, nil, fmt.Errorf("Error generating a random block ID for source file %q: %s", file.Name(), err)
		}

		sectionSize := blockSize
		remainder := info.Size() - i
		if remainder < blockSize {
			sectionSize = remainder
		}

		block := storage.Block{
			ID:     base64.StdEncoding.EncodeToString(entropy),
			Status: storage.BlockStatusUncommitted,
		}

		blockList = append(blockList, block)

		parts = append(parts, resourceArmStorageBlobBlock{
			id:      block.ID,
			section: io.NewSectionReader(file, i, sectionSize),
		})
	}

	return blockList, parts, nil
}

type resourceArmStorageBlobBlockUploadContext struct {
	client    *storage.BlobStorageClient
	container string
	name      string
	source    string
	attempts  int
	blocks    chan resourceArmStorageBlobBlock
	errors    chan error
	wg        *sync.WaitGroup
}

func resourceArmStorageBlobBlockUploadWorker(ctx resourceArmStorageBlobBlockUploadContext) {
	for block := range ctx.blocks {
		buffer := make([]byte, block.section.Size())

		_, err := block.section.Read(buffer)
		if err != nil {
			ctx.errors <- fmt.Errorf("Error reading source file %q: %s", ctx.source, err)
			ctx.wg.Done()
			continue
		}

		for i := 0; i < ctx.attempts; i++ {
			container := ctx.client.GetContainerReference(ctx.container)
			blob := container.GetBlobReference(ctx.name)
			options := &storage.PutBlockOptions{}
			if err = blob.PutBlock(block.id, buffer, options); err == nil {
				break
			}
		}
		if err != nil {
			ctx.errors <- fmt.Errorf("Error uploading block %q for source file %q: %s", block.id, ctx.source, err)
			ctx.wg.Done()
			continue
		}

		ctx.wg.Done()
	}
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

	blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
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

	blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
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

	blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, *resourceGroup, id.storageAccountName)
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

type storageBlobId struct {
	storageAccountName string
	containerName      string
	blobName           string
}

func parseStorageBlobID(input string, environment azauto.Environment) (*storageBlobId, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("Error parsing %q as URI: %+v", input, err)
	}

	// trim the leading `/`
	segments := strings.Split(strings.TrimPrefix(uri.Path, "/"), "/")
	if len(segments) < 2 {
		return nil, fmt.Errorf("Expected number of segments in the path to be < 2 but got %d", len(segments))
	}

	storageAccountName := strings.Replace(uri.Host, fmt.Sprintf(".blob.%s", environment.StorageEndpointSuffix), "", 1)
	containerName := segments[0]
	blobName := strings.TrimPrefix(uri.Path, fmt.Sprintf("/%s/", containerName))

	id := storageBlobId{
		storageAccountName: storageAccountName,
		containerName:      containerName,
		blobName:           blobName,
	}
	return &id, nil
}

func determineResourceGroupForStorageAccount(accountName string, client *ArmClient) (*string, error) {
	storageClient := client.storageServiceClient
	ctx := client.StopContext

	// first locate which resource group the storage account is in
	groupsResp, err := storageClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error loading the Resource Groups for Storage Account %q: %+v", accountName, err)
	}

	if groups := groupsResp.Value; groups != nil {
		for _, group := range *groups {
			if group.Name != nil && *group.Name == accountName {
				groupId, err := azure.ParseAzureResourceID(*group.ID)
				if err != nil {
					return nil, err
				}

				return &groupId.ResourceGroup, nil
			}
		}
	}

	return nil, nil
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
