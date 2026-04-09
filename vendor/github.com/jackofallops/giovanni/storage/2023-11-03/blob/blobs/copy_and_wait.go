package blobs

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

// CopyAndWait copies a blob to a destination within the storage account and waits for it to finish copying.
func (c Client) CopyAndWait(ctx context.Context, containerName, blobName string, input CopyInput) error {
	if _, err := c.Copy(ctx, containerName, blobName, input); err != nil {
		return fmt.Errorf("error copying: %s", err)
	}

	getInput := GetPropertiesInput{
		LeaseID: input.LeaseID,
	}

	pollerType := NewCopyAndWaitPoller(&c, containerName, blobName, getInput)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for file to copy: %+v", err)
	}

	return nil
}
