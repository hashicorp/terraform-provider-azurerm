package files

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

// CopyAndWait is a convenience method which doesn't exist in the API, which copies the file and then waits for the copy to complete
func (c Client) CopyAndWait(ctx context.Context, shareName, path, fileName string, input CopyInput) (resp CopyResponse, err error) {
	copy, e := c.Copy(ctx, shareName, path, fileName, input)
	if err != nil {
		resp.HttpResponse = copy.HttpResponse
		err = fmt.Errorf("error copying: %s", e)
		return
	}

	resp.CopyID = copy.CopyID

	pollerType := NewCopyAndWaitPoller(&c, shareName, path, fileName)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return resp, fmt.Errorf("waiting for file to copy: %+v", err)
	}

	return
}
