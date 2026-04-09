package files

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane/storage"
)

// Client is the base client for File Storage Shares.
type Client struct {
	Client *storage.Client
}

func NewWithBaseUri(baseUri string) (*Client, error) {
	baseClient, err := storage.NewStorageClient(baseUri, componentName, apiVersion)
	if err != nil {
		return nil, fmt.Errorf("building base client: %+v", err)
	}

	baseClient.Client.AuthorizeRequest = func(ctx context.Context, req *http.Request, authorizer auth.Authorizer) error {
		if err := auth.SetAuthHeader(ctx, req, authorizer); err != nil {
			return fmt.Errorf("authorizing request: %+v", err)
		}

		// Only set this header if OAuth is being used (i.e. not shared key authentication)
		if _, ok := authorizer.(*auth.SharedKeyAuthorizer); !ok {
			req.Header.Set("x-ms-file-request-intent", "backup")
		}

		return nil
	}

	return &Client{
		Client: baseClient,
	}, nil
}
