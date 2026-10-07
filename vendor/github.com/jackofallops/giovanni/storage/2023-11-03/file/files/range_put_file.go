package files

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sync"
)

// PutFile is a helper method which takes a file, and automatically chunks it up, rather than having to do this yourself
func (c Client) PutFile(ctx context.Context, shareName, path, fileName string, file *os.File, parallelism int) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error loading file info: %s", err)
	}
	if fileInfo.Size() == 0 {
		return fmt.Errorf("file is empty which is not supported")
	}

	fileSize := fileInfo.Size()
	chunkSize := 4 * 1024 * 1024 // 4MB
	if chunkSize > int(fileSize) {
		chunkSize = int(fileSize)
	}
	chunks := int(math.Ceil(float64(fileSize) / float64(chunkSize*1.0)))

	workerCount := parallelism
	if workerCount > chunks {
		workerCount = chunks
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(workerCount)

	jobs := make(chan int, workerCount)
	errors := make(chan error, chunkSize)

	for i := 0; i < workerCount; i++ {
		go func() {
			for i := range jobs {
				log.Printf("[DEBUG] Chunk %d of %d", i+1, chunks)

				uci := uploadChunkInput{
					thisChunk: i,
					chunkSize: chunkSize,
					fileSize:  fileSize,
				}

				_, err := c.uploadChunk(ctx, shareName, path, fileName, uci, file)
				if err != nil {
					errors <- err
				}
			}
			waitGroup.Done()
		}()
	}

	for i := 0; i < chunks; i++ {
		jobs <- i
	}
	close(jobs)
	waitGroup.Wait()

	// TODO: we should switch to hashicorp/multi-error here
	if len(errors) > 0 {
		return fmt.Errorf("uploading file: %s", <-errors)
	}

	return nil
}

type uploadChunkInput struct {
	thisChunk int
	chunkSize int
	fileSize  int64
}

func (c Client) uploadChunk(ctx context.Context, shareName, path, fileName string, input uploadChunkInput, file *os.File) (result PutRangeResponse, err error) {
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
			return result, fmt.Errorf("reading bytes: %s", err)
		}
	}

	putBytesInput := PutByteRangeInput{
		StartBytes: startBytes,
		EndBytes:   endBytes,
		Content:    bytes,
	}
	result, err = c.PutByteRange(ctx, shareName, path, fileName, putBytesInput)
	if err != nil {
		return result, fmt.Errorf("putting bytes: %s", err)
	}

	return
}
