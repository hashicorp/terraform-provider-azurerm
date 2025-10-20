package files

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

// CopyAndWait is a convenience method which doesn't exist in the API, which copies the file and then waits for the copy to complete
func (c Client) CopyAndWait(ctx context.Context, shareName, path, fileName string, input CopyInput) (result CopyResponse, err error) {
	fileCopy, e := c.Copy(ctx, shareName, path, fileName, input)
	if err != nil {
		result.HttpResponse = fileCopy.HttpResponse
		err = fmt.Errorf("copying: %s", e)
		return
	}

	result.CopyID = fileCopy.CopyID

	pollerType := NewCopyAndWaitPoller(&c, shareName, path, fileName)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err = poller.PollUntilDone(ctx); err != nil {
		return result, fmt.Errorf("waiting for file to copy: %+v", err)
	}

	return
}
