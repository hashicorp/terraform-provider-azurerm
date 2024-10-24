package custompollers

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/fileshares"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = storageShareCreatePoller{}

type storageShareCreatePoller struct {
	id      fileshares.ShareId
	client  *fileshares.FileSharesClient
	payload fileshares.FileShare
}

func NewStorageShareCreatePoller(client *fileshares.FileSharesClient, id fileshares.ShareId, payload fileshares.FileShare) *storageShareCreatePoller {
	return &storageShareCreatePoller{
		id:      id,
		client:  client,
		payload: payload,
	}
}

func (p storageShareCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	// Note - Whilst this is an antipattern for the Provider, the API provides no way currently to poll for deletion
	// to ensure it's removed. To support rapid delete then re-creation we check for 409's that indicate the resource
	// is still being removed.
	resp, err := p.client.Create(ctx, p.id, p.payload, fileshares.DefaultCreateOperationOptions())
	if err != nil {
		if response.WasConflict(resp.HttpResponse) {
			return &pollers.PollResult{
				PollInterval: 5 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}

		return &pollers.PollResult{
			HttpResponse: nil,
			PollInterval: 5 * time.Second,
			Status:       pollers.PollingStatusFailed,
		}, err
	}

	return &pollers.PollResult{
		PollInterval: 5,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
