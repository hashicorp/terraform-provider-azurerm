package files

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &copyAndWaitPoller{}

func NewCopyAndWaitPoller(client *Client, shareName, path, fileName string) *copyAndWaitPoller {
	return &copyAndWaitPoller{
		client:    client,
		shareName: shareName,
		path:      path,
		fileName:  fileName,
	}
}

type copyAndWaitPoller struct {
	client    *Client
	shareName string
	path      string
	fileName  string
}

func (p *copyAndWaitPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	props, err := p.client.GetProperties(ctx, p.shareName, p.path, p.fileName)
	if err != nil {
		return nil, fmt.Errorf("retrieving copy (shareName: %s path: %s fileName: %s) : %+v", p.shareName, p.path, p.fileName, err)
	}

	if strings.EqualFold(props.CopyStatus, "success") {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusSucceeded,
			PollInterval: 10 * time.Second,
		}, nil
	}

	// Processing
	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}, nil
}
