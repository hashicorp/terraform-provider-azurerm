package blobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &copyAndWaitPoller{}

func NewCopyAndWaitPoller(client *Client, containerName, blobName string, getPropertiesInput GetPropertiesInput) *copyAndWaitPoller {
	return &copyAndWaitPoller{
		client:             client,
		containerName:      containerName,
		blobName:           blobName,
		getPropertiesInput: getPropertiesInput,
	}
}

type copyAndWaitPoller struct {
	client             *Client
	containerName      string
	blobName           string
	getPropertiesInput GetPropertiesInput
}

func (p *copyAndWaitPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	props, err := p.client.GetProperties(ctx, p.containerName, p.blobName, p.getPropertiesInput)
	if err != nil {
		return nil, fmt.Errorf("retrieving properties (container: %s blob: %s) : %+v", p.containerName, p.blobName, err)
	}

	if strings.EqualFold(string(props.CopyStatus), string(Success)) {
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
