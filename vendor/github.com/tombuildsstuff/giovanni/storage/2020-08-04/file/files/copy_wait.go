package files

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
)

type CopyAndWaitResult struct {
	autorest.Response

	CopyID string
}

const DefaultCopyPollDuration = 15 * time.Second

// CopyAndWait is a convenience method which doesn't exist in the API, which copies the file and then waits for the copy to complete
func (client Client) CopyAndWait(ctx context.Context, accountName, shareName, path, fileName string, input CopyInput, pollDuration time.Duration) (result CopyResult, err error) {
	copy, e := client.Copy(ctx, accountName, shareName, path, fileName, input)
	if err != nil {
		result.Response = copy.Response
		err = fmt.Errorf("Error copying: %s", e)
		return
	}

	result.CopyID = copy.CopyID

	// since the API doesn't return a LRO, this is a hack which also polls every 10s, but should be sufficient
	for true {
		props, e := client.GetProperties(ctx, accountName, shareName, path, fileName)
		if e != nil {
			result.Response = copy.Response
			err = fmt.Errorf("Error waiting for copy: %s", e)
			return
		}

		switch strings.ToLower(props.CopyStatus) {
		case "pending":
			time.Sleep(pollDuration)
			continue

		case "success":
			return

		default:
			err = fmt.Errorf("Unexpected CopyState %q", e)
			return
		}
	}

	return
}
