package files

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sync"

	"github.com/Azure/go-autorest/autorest"
)

// PutFile is a helper method which takes a file, and automatically chunks it up, rather than having to do this yourself
func (client Client) PutFile(ctx context.Context, accountName, shareName, path, fileName string, file *os.File, parallelism int) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Error loading file info: %s", err)
	}

	fileSize := fileInfo.Size()
	chunkSize := 4 * 1024 * 1024 // 4MB
	if chunkSize > int(fileSize) {
		chunkSize = int(fileSize)
	}
	chunks := int(math.Ceil(float64(fileSize) / float64(chunkSize*1.0)))

	workerCount := parallelism * runtime.NumCPU()
	if workerCount > chunks {
		workerCount = chunks
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(workerCount)
	errors := make(chan error, chunkSize)

	for i := 0; i < chunks; i++ {
		go func(i int) {
			log.Printf("[DEBUG] Chunk %d of %d", i+1, chunks)

			uci := uploadChunkInput{
				thisChunk: i,
				chunkSize: chunkSize,
				fileSize:  fileSize,
			}

			_, err := client.uploadChunk(ctx, accountName, shareName, path, fileName, uci, file)
			if err != nil {
				errors <- err
				waitGroup.Done()
				return
			}

			waitGroup.Done()
			return
		}(i)
	}
	waitGroup.Wait()

	// TODO: we should switch to hashicorp/multi-error here
	if len(errors) > 0 {
		return fmt.Errorf("Error uploading file: %s", <-errors)
	}

	return nil
}

type uploadChunkInput struct {
	thisChunk int
	chunkSize int
	fileSize  int64
}

func (client Client) uploadChunk(ctx context.Context, accountName, shareName, path, fileName string, input uploadChunkInput, file *os.File) (result autorest.Response, err error) {
	startBytes := int64(input.chunkSize * input.thisChunk)
	endBytes := startBytes + int64(input.chunkSize)

	// the last size may exceed the size of the file
	remaining := input.fileSize - startBytes
	if int64(input.chunkSize) > remaining {
		endBytes = startBytes + remaining
	}

	bytesToRead := int(endBytes) - int(startBytes)
	bytes := make([]byte, bytesToRead)

	_, err = file.ReadAt(bytes, startBytes)
	if err != nil {
		if err != io.EOF {
			return result, fmt.Errorf("Error reading bytes: %s", err)
		}
	}

	putBytesInput := PutByteRangeInput{
		StartBytes: startBytes,
		EndBytes:   endBytes,
		Content:    bytes,
	}
	result, err = client.PutByteRange(ctx, accountName, shareName, path, fileName, putBytesInput)
	if err != nil {
		return result, fmt.Errorf("Error putting bytes: %s", err)
	}

	return
}
