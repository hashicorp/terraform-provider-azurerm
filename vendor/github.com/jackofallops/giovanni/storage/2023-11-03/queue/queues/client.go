package queues

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane/storage"
)

// Client is the base client for Queue Storage Shares.

// Client is the base client for Messages.
type Client struct {
	Client *storage.Client
}

func NewWithBaseUri(baseUri string) (*Client, error) {
	baseClient, err := storage.NewStorageClient(baseUri, componentName, apiVersion)
	if err != nil {
		return nil, fmt.Errorf("building base client: %+v", err)
	}
	return &Client{
		Client: baseClient,
	}, nil
}
