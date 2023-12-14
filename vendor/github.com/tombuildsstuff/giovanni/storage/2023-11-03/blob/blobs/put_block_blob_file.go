package blobs

import (
	"context"
	"fmt"
	"io"
	"os"
)

// PutBlockBlobFromFile is a helper method which takes a file, and automatically chunks it up, rather than having to do this yourself
func (c Client) PutBlockBlobFromFile(ctx context.Context, containerName, blobName string, file *os.File, input PutBlockBlobInput) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error loading file info: %s", err)
	}

	fileSize := fileInfo.Size()
	bytes := make([]byte, fileSize)

	_, err = file.ReadAt(bytes, 0)
	if err != nil {
		if err != io.EOF {
			return fmt.Errorf("Error reading bytes: %s", err)
		}
	}

	input.Content = &bytes

	if _, err = c.PutBlockBlob(ctx, containerName, blobName, input); err != nil {
		return fmt.Errorf("error putting bytes: %s", err)
	}

	return nil
}
