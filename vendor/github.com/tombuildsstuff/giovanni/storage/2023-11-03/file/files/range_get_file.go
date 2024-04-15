package files

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"runtime"
	"sync"
)

type GetFileInput struct {
	Parallelism int
}

type GetFileResponse struct {
	HttpResponse *http.Response
	OutputBytes  *[]byte
}

// GetFile is a helper method to download a file by chunking it automatically
func (c Client) GetFile(ctx context.Context, shareName, path, fileName string, input GetFileInput) (result GetFileResponse, err error) {

	// first look up the file and check out how many bytes it is
	file, e := c.GetProperties(ctx, shareName, path, fileName)
	if err != nil {
		result.HttpResponse = file.HttpResponse
		err = e
		return
	}

	if file.ContentLength == nil {
		err = fmt.Errorf("Content-Length was nil")
		return
	}

	result.HttpResponse = file.HttpResponse
	length := *file.ContentLength
	chunkSize := int64(4 * 1024 * 1024) // 4MB

	if chunkSize > length {
		chunkSize = length
	}

	// then split that up into chunks and retrieve it into the 'results' set
	chunks := int(math.Ceil(float64(length) / float64(chunkSize)))
	workerCount := input.Parallelism * runtime.NumCPU()
	if workerCount > chunks {
		workerCount = chunks
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(workerCount)

	results := make([]*downloadFileChunkResult, chunks)
	errors := make(chan error, chunkSize)

	for i := 0; i < chunks; i++ {
		go func(i int) {
			log.Printf("[DEBUG] Downloading Chunk %d of %d", i+1, chunks)

			dfci := downloadFileChunkInput{
				thisChunk: i,
				chunkSize: chunkSize,
				fileSize:  length,
			}

			result, err := c.downloadFileChunk(ctx, shareName, path, fileName, dfci)
			if err != nil {
				errors <- err
				waitGroup.Done()
				return
			}

			// if there's no error, we should have bytes, so this is safe
			results[i] = result

			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()

	// TODO: we should switch to hashicorp/multi-error here
	if len(errors) > 0 {
		err = fmt.Errorf("Error downloading file: %s", <-errors)
		return
	}

	// then finally put it all together, in order and return it
	output := make([]byte, length)
	for _, v := range results {
		copy(output[v.startBytes:v.endBytes], *v.bytes)
	}

	if result.OutputBytes == nil {
		result.OutputBytes = &[]byte{}
	}
	*result.OutputBytes = output
	return
}

type downloadFileChunkInput struct {
	thisChunk int
	chunkSize int64
	fileSize  int64
}

type downloadFileChunkResult struct {
	startBytes int64
	endBytes   int64
	bytes      *[]byte
}

func (c Client) downloadFileChunk(ctx context.Context, shareName, path, fileName string, input downloadFileChunkInput) (*downloadFileChunkResult, error) {
	startBytes := input.chunkSize * int64(input.thisChunk)
	endBytes := startBytes + input.chunkSize

	// the last chunk may exceed the size of the file
	remaining := input.fileSize - startBytes
	if input.chunkSize > remaining {
		endBytes = startBytes + remaining
	}

	getInput := GetByteRangeInput{
		StartBytes: startBytes,
		EndBytes:   endBytes,
	}
	result, err := c.GetByteRange(ctx, shareName, path, fileName, getInput)
	if err != nil {
		return nil, fmt.Errorf("error putting bytes: %s", err)
	}

	output := downloadFileChunkResult{
		startBytes: startBytes,
		endBytes:   endBytes,
		bytes:      result.Contents,
	}
	return &output, nil
}
