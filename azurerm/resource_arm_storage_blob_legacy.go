package azurerm

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
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
	container := sbu.legacyClient.GetContainerReference(sbu.containerName)
	blob := container.GetBlobReference(sbu.blobName)

	if sbu.sourceUri != "" {
		return sbu.copy(ctx)
	}

	switch strings.ToLower(sbu.blobType) {
	// TODO: new feature for 'append' blobs?
	case "block":
		options := &legacy.PutBlobOptions{}
		if err := blob.CreateBlockBlob(options); err != nil {
			return fmt.Errorf("Error creating storage blob on Azure: %s", err)
		}

		if sbu.source != "" {
			if err := resourceArmStorageBlobBlockUploadFromSource(sbu.containerName, sbu.blobName, sbu.source, sbu.contentType, sbu.legacyClient, sbu.parallelism, sbu.attempts); err != nil {
				return fmt.Errorf("Error creating storage blob on Azure: %s", err)
			}
		}
	case "page":
		if sbu.source != "" {
			if err := resourceArmStorageBlobPageUploadFromSource(sbu.containerName, sbu.blobName, sbu.source, sbu.contentType, sbu.legacyClient, sbu.parallelism, sbu.attempts); err != nil {
				return fmt.Errorf("Error creating storage blob on Azure: %s", err)
			}
		} else {
			size := int64(sbu.size)
			options := &legacy.PutBlobOptions{}

			blob.Properties.ContentLength = size
			blob.Properties.ContentType = sbu.contentType
			if err := blob.PutPageBlob(options); err != nil {
				return fmt.Errorf("Error creating storage blob on Azure: %s", err)
			}
		}
	}

	return nil
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

// TODO: remove below here

type resourceArmStorageBlobPage struct {
	offset  int64
	section *io.SectionReader
}

func resourceArmStorageBlobPageUploadFromSource(container, name, source, contentType string, client *legacy.BlobStorageClient, parallelism, attempts int) error {
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

	options := &legacy.PutBlobOptions{}
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
	client    *legacy.BlobStorageClient
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

func resourceArmStorageBlobBlockUploadFromSource(container, name, source, contentType string, client *legacy.BlobStorageClient, parallelism, attempts int) error {
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
	options := &legacy.PutBlockListOptions{}
	err = blobReference.PutBlockList(blockList, options)
	if err != nil {
		return fmt.Errorf("Error updating block list for source file %q: %s", source, err)
	}

	return nil
}

func resourceArmStorageBlobBlockSplit(file *os.File) ([]legacy.Block, []resourceArmStorageBlobBlock, error) {
	const (
		idSize          = 64
		blockSize int64 = 4 * 1024 * 1024
	)
	var parts []resourceArmStorageBlobBlock
	var blockList []legacy.Block

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

		block := legacy.Block{
			ID:     base64.StdEncoding.EncodeToString(entropy),
			Status: legacy.BlockStatusUncommitted,
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
	client    *legacy.BlobStorageClient
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
			options := &legacy.PutBlockOptions{}
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
