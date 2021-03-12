package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/blobs"
)

const pollingInterval = time.Second * 15

type BlobUpload struct {
	Client *blobs.Client

	AccountName   string
	BlobName      string
	ContainerName string

	BlobType      string
	ContentType   string
	ContentMD5    string
	MetaData      map[string]string
	Parallelism   int
	Size          int
	Source        string
	SourceContent string
	SourceUri     string
}

func (sbu BlobUpload) Create(ctx context.Context) error {
	blobType := strings.ToLower(sbu.BlobType)

	if blobType == "append" {
		if sbu.Source != "" || sbu.SourceContent != "" || sbu.SourceUri != "" {
			return fmt.Errorf("A source cannot be specified for an Append blob")
		}

		if sbu.ContentMD5 != "" {
			return fmt.Errorf("`content_md5` cannot be specified for an Append blob")
		}

		return sbu.createEmptyAppendBlob(ctx)
	}

	if blobType == "block" {
		if sbu.SourceUri != "" {
			return sbu.copy(ctx)
		}

		if sbu.SourceContent != "" {
			return sbu.uploadBlockBlobFromContent(ctx)
		}
		if sbu.Source != "" {
			return sbu.uploadBlockBlob(ctx)
		}

		return sbu.createEmptyBlockBlob(ctx)
	}

	if blobType == "page" {
		if sbu.ContentMD5 != "" {
			return fmt.Errorf("`content_md5` cannot be specified for a Page blob")
		}
		if sbu.SourceUri != "" {
			return sbu.copy(ctx)
		}
		if sbu.SourceContent != "" {
			return fmt.Errorf("`source_content` cannot be specified for a Page blob")
		}
		if sbu.Source != "" {
			return sbu.uploadPageBlob(ctx)
		}

		return sbu.createEmptyPageBlob(ctx)
	}

	return fmt.Errorf("Unsupported Blob Type: %q", blobType)
}

func (sbu BlobUpload) copy(ctx context.Context) error {
	input := blobs.CopyInput{
		CopySource: sbu.SourceUri,
		MetaData:   sbu.MetaData,
	}
	if err := sbu.Client.CopyAndWait(ctx, sbu.AccountName, sbu.ContainerName, sbu.BlobName, input, pollingInterval); err != nil {
		return fmt.Errorf("Error copy/waiting: %s", err)
	}

	return nil
}

func (sbu BlobUpload) createEmptyAppendBlob(ctx context.Context) error {
	input := blobs.PutAppendBlobInput{
		ContentType: utils.String(sbu.ContentType),
		MetaData:    sbu.MetaData,
	}
	if _, err := sbu.Client.PutAppendBlob(ctx, sbu.AccountName, sbu.ContainerName, sbu.BlobName, input); err != nil {
		return fmt.Errorf("Error PutAppendBlob: %s", err)
	}

	return nil
}

func (sbu BlobUpload) createEmptyBlockBlob(ctx context.Context) error {
	if sbu.ContentMD5 != "" {
		return fmt.Errorf("`content_md5` cannot be specified for empty Block blobs")
	}

	input := blobs.PutBlockBlobInput{
		ContentType: utils.String(sbu.ContentType),
		MetaData:    sbu.MetaData,
	}
	if _, err := sbu.Client.PutBlockBlob(ctx, sbu.AccountName, sbu.ContainerName, sbu.BlobName, input); err != nil {
		return fmt.Errorf("Error PutBlockBlob: %s", err)
	}

	return nil
}

func (sbu BlobUpload) uploadBlockBlobFromContent(ctx context.Context) error {
	tmpFile, err := os.CreateTemp(os.TempDir(), "upload-")
	if err != nil {
		return fmt.Errorf("Error creating temporary file: %s", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write([]byte(sbu.SourceContent)); err != nil {
		return fmt.Errorf("Error writing Source Content to Temp File: %s", err)
	}
	defer tmpFile.Close()

	sbu.Source = tmpFile.Name()
	return sbu.uploadBlockBlob(ctx)
}

func (sbu BlobUpload) uploadBlockBlob(ctx context.Context) error {
	file, err := os.Open(sbu.Source)
	if err != nil {
		return fmt.Errorf("Error opening: %s", err)
	}
	defer file.Close()

	input := blobs.PutBlockBlobInput{
		ContentType: utils.String(sbu.ContentType),
		MetaData:    sbu.MetaData,
	}
	if sbu.ContentMD5 != "" {
		input.ContentMD5 = utils.String(sbu.ContentMD5)
	}
	if err := sbu.Client.PutBlockBlobFromFile(ctx, sbu.AccountName, sbu.ContainerName, sbu.BlobName, file, input); err != nil {
		return fmt.Errorf("Error PutBlockBlobFromFile: %s", err)
	}

	return nil
}

func (sbu BlobUpload) createEmptyPageBlob(ctx context.Context) error {
	if sbu.Size == 0 {
		return fmt.Errorf("`size` cannot be zero for a page blob")
	}

	input := blobs.PutPageBlobInput{
		BlobContentLengthBytes: int64(sbu.Size),
		ContentType:            utils.String(sbu.ContentType),
		MetaData:               sbu.MetaData,
	}
	if _, err := sbu.Client.PutPageBlob(ctx, sbu.AccountName, sbu.ContainerName, sbu.BlobName, input); err != nil {
		return fmt.Errorf("Error PutPageBlob: %s", err)
	}

	return nil
}

func (sbu BlobUpload) uploadPageBlob(ctx context.Context) error {
	if sbu.Size != 0 {
		return fmt.Errorf("`size` cannot be set for an uploaded page blob")
	}

	// determine the details about the file
	file, err := os.Open(sbu.Source)
	if err != nil {
		return fmt.Errorf("Error opening source file for upload %q: %s", sbu.Source, err)
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
		BlobContentLengthBytes: fileSize,
		ContentType:            utils.String(sbu.ContentType),
		MetaData:               sbu.MetaData,
	}
	if _, err := sbu.Client.PutPageBlob(ctx, sbu.AccountName, sbu.ContainerName, sbu.BlobName, input); err != nil {
		return fmt.Errorf("Error PutPageBlob: %s", err)
	}

	if err := sbu.pageUploadFromSource(ctx, file, fileSize); err != nil {
		return fmt.Errorf("Error creating storage blob on Azure: %s", err)
	}

	return nil
}

// TODO: move below here into Giovanni

type storageBlobPage struct {
	offset  int64
	section *io.SectionReader
}

func (sbu BlobUpload) pageUploadFromSource(ctx context.Context, file io.ReaderAt, fileSize int64) error {
	workerCount := sbu.Parallelism * runtime.NumCPU()

	// first we chunk the file and assign them to 'pages'
	pageList, err := sbu.storageBlobPageSplit(file, fileSize)
	if err != nil {
		return fmt.Errorf("Error splitting source file %q into pages: %s", sbu.Source, err)
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
		go sbu.blobPageUploadWorker(ctx, blobPageUploadContext{
			blobSize: fileSize,
			pages:    pages,
			errors:   errors,
			wg:       wg,
		})
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("Error while uploading source file %q: %s", sbu.Source, <-errors)
	}

	return nil
}

const (
	minPageSize int64 = 4 * 1024

	// TODO: investigate whether this can be bumped to 100MB with the new API
	maxPageSize int64 = 4 * 1024 * 1024
)

func (sbu BlobUpload) storageBlobPageSplit(file io.ReaderAt, fileSize int64) ([]storageBlobPage, error) {
	// whilst the file Size can be any arbitrary Size, it must be uploaded in fixed-Size pages
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
		if _, err := file.ReadAt(pageBuf, i); err != nil && err != io.EOF {
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
}

func (sbu BlobUpload) blobPageUploadWorker(ctx context.Context, uploadCtx blobPageUploadContext) {
	for page := range uploadCtx.pages {
		start := page.offset
		end := page.offset + page.section.Size() - 1
		if end > uploadCtx.blobSize-1 {
			end = uploadCtx.blobSize - 1
		}
		size := end - start + 1

		chunk := make([]byte, size)
		if _, err := page.section.Read(chunk); err != nil && err != io.EOF {
			uploadCtx.errors <- fmt.Errorf("Error reading source file %q at offset %d: %s", sbu.Source, page.offset, err)
			uploadCtx.wg.Done()
			continue
		}

		input := blobs.PutPageUpdateInput{
			StartByte: start,
			EndByte:   end,
			Content:   chunk,
		}

		if _, err := sbu.Client.PutPageUpdate(ctx, sbu.AccountName, sbu.ContainerName, sbu.BlobName, input); err != nil {
			uploadCtx.errors <- fmt.Errorf("Error writing page at offset %d for file %q: %s", page.offset, sbu.Source, err)
			uploadCtx.wg.Done()
			continue
		}

		uploadCtx.wg.Done()
	}
}

func convertHexToBase64Encoding(str string) (string, error) {
	data, err := hex.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("converting %q from Hex to Base64 Encoding: %+v", str, err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func convertBase64ToHexEncoding(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("converting %q from Base64 to Hex Encoding: %+v", str, err)
	}

	return hex.EncodeToString(data), nil
}
