package azurerm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	legacy "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/blobs"
)

const pollingInterval = time.Second * 15

type StorageBlobUpload struct {
	accountName   string
	containerName string
	blobName      string

	attempts    int
	blobType    string
	contentType string
	metaData    map[string]string
	parallelism int
	size        int
	source      string
	sourceUri   string

	client       *blobs.Client
	legacyClient *legacy.BlobStorageClient
}

func (sbu StorageBlobUpload) Create(ctx context.Context) error {
	if sbu.sourceUri != "" {
		return sbu.copy(ctx)
	}

	blobType := strings.ToLower(sbu.blobType)

	// TODO: new feature for 'append' blobs?

	if blobType == "block" {
		if sbu.source != "" {
			return sbu.uploadBlockBlob(ctx)
		}

		return sbu.createEmptyBlockBlob(ctx)
	}

	if blobType == "page" {
		if sbu.source != "" {
			return sbu.uploadPageBlob(ctx)
		}

		return sbu.createEmptyPageBlob(ctx)
	}

	return fmt.Errorf("Unsupported Blob Type: %q", blobType)
}

func (sbu StorageBlobUpload) copy(ctx context.Context) error {
	input := blobs.CopyInput{
		CopySource: sbu.sourceUri,
		MetaData:   sbu.metaData,
	}
	if err := sbu.client.CopyAndWait(ctx, sbu.accountName, sbu.containerName, sbu.blobName, input, pollingInterval); err != nil {
		return fmt.Errorf("Error copy/waiting: %s", err)
	}

	return nil
}

func (sbu StorageBlobUpload) createEmptyBlockBlob(ctx context.Context) error {
	input := blobs.PutBlockBlobInput{
		ContentType: utils.String(sbu.contentType),
		MetaData:    sbu.metaData,
	}
	if _, err := sbu.client.PutBlockBlob(ctx, sbu.accountName, sbu.containerName, sbu.blobName, input); err != nil {
		return fmt.Errorf("Error PutBlockBlob: %s", err)
	}

	return nil
}

func (sbu StorageBlobUpload) uploadBlockBlob(ctx context.Context) error {
	file, err := os.Open(sbu.source)
	if err != nil {
		return fmt.Errorf("Error opening: %s", err)
	}
	defer file.Close()

	input := blobs.PutBlockBlobInput{
		ContentType: utils.String(sbu.contentType),
		MetaData:    sbu.metaData,
	}

	if err := sbu.client.PutBlockBlobFromFile(ctx, sbu.accountName, sbu.containerName, sbu.blobName, file, input); err != nil {
		return fmt.Errorf("Error PutBlockBlobFromFile: %s", err)
	}

	return nil
}

func (sbu StorageBlobUpload) createEmptyPageBlob(ctx context.Context) error {
	if sbu.size == 0 {
		return fmt.Errorf("`size` cannot be zero for a page blob")
	}

	input := blobs.PutPageBlobInput{
		BlobContentLengthBytes: int64(sbu.size),
		ContentType:            utils.String(sbu.contentType),
		MetaData:               sbu.metaData,
	}
	// TODO: access tiers?
	if _, err := sbu.client.PutPageBlob(ctx, sbu.accountName, sbu.containerName, sbu.blobName, input); err != nil {
		return fmt.Errorf("Error PutPageBlob: %s", err)
	}

	return nil
}

func (sbu StorageBlobUpload) uploadPageBlob(ctx context.Context) error {
	if sbu.size != 0 {
		// the user shouldn't need to specify this since we infer it
	}

	// determine the details about the file
	file, err := os.Open(sbu.source)
	if err != nil {
		return fmt.Errorf("Error opening source file for upload %q: %s", sbu.source, err)
	}
	defer file.Close()

	// TODO: all of this ultimately can be moved into Giovanni

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Could not stat file %q: %s", file.Name(), err)
	}

	fileSize := info.Size()

	// first let's create a file of the specified file size
	input := blobs.PutPageBlobInput{
		BlobContentLengthBytes: int64(fileSize),
		ContentType:            utils.String(sbu.contentType),
		MetaData:               sbu.metaData,
	}
	// TODO: access tiers?
	if _, err := sbu.client.PutPageBlob(ctx, sbu.accountName, sbu.containerName, sbu.blobName, input); err != nil {
		return fmt.Errorf("Error PutPageBlob: %s", err)
	}

	if err := sbu.pageUploadFromSource(ctx, file, fileSize); err != nil {
		return fmt.Errorf("Error creating storage blob on Azure: %s", err)
	}

	return nil
}

// TODO: remove below here

type storageBlobPage struct {
	offset  int64
	section *io.SectionReader
}

func (sbu StorageBlobUpload) pageUploadFromSource(ctx context.Context, file *os.File, fileSize int64) error {
	workerCount := sbu.parallelism * runtime.NumCPU()

	// first we chunk the file and assign them to 'pages'
	pageList, err := sbu.resourceArmStorageBlobPageSplit(file, fileSize)
	if err != nil {
		return fmt.Errorf("Error splitting source file %q into pages: %s", sbu.source, err)
	}

	// finally we upload the contents of said file
	pages := make(chan storageBlobPage, len(pageList))
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
		go sbu.blobPageUploadWorker(blobPageUploadContext{
			blobSize: fileSize,
			pages:    pages,
			errors:   errors,
			wg:       wg,
			attempts: sbu.attempts,
		})
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("Error while uploading source file %q: %s", sbu.source, <-errors)
	}

	return nil
}

const (
	minPageSize int64 = 4 * 1024
	maxPageSize int64 = 4 * 1024 * 1024
)

func (sbu StorageBlobUpload) resourceArmStorageBlobPageSplit(file *os.File, fileSize int64) ([]storageBlobPage, error) {
	// whilst the file size can be any arbitary size, it must be uploaded in fixed-size pages
	blobSize := fileSize
	if fileSize%minPageSize != 0 {
		blobSize = fileSize + (minPageSize - (fileSize % minPageSize))
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
		_, err := file.ReadAt(pageBuf, i)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("Could not read chunk at %d: %s", i, err)
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

	var pages []storageBlobPage
	for _, nonEmptyRange := range nonEmptyRanges {
		pages = append(pages, storageBlobPage{
			offset:  nonEmptyRange.offset,
			section: io.NewSectionReader(file, nonEmptyRange.offset, nonEmptyRange.length),
		})
	}

	return pages, nil
}

type blobPageUploadContext struct {
	blobSize int64
	pages    chan storageBlobPage
	errors   chan error
	wg       *sync.WaitGroup
	attempts int
}

func (sbu StorageBlobUpload) blobPageUploadWorker(ctx blobPageUploadContext) {
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
			ctx.errors <- fmt.Errorf("Error reading source file %q at offset %d: %s", sbu.source, page.offset, err)
			ctx.wg.Done()
			continue
		}

		for x := 0; x < ctx.attempts; x++ {
			container := sbu.legacyClient.GetContainerReference(sbu.containerName)
			blob := container.GetBlobReference(sbu.blobName)
			blobRange := legacy.BlobRange{
				Start: uint64(start),
				End:   uint64(end),
			}
			options := &legacy.PutPageOptions{}
			reader := bytes.NewReader(chunk)
			err = blob.WriteRange(blobRange, reader, options)
			if err == nil {
				break
			}
		}
		if err != nil {
			ctx.errors <- fmt.Errorf("Error writing page at offset %d for file %q: %s", page.offset, sbu.source, err)
			ctx.wg.Done()
			continue
		}

		ctx.wg.Done()
	}
}
