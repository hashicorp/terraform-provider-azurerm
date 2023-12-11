package paths

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Client is the base client for Data Lake Storage Path
type Client struct {
	autorest.Client
	BaseURI string
}

// New creates an instance of the Data Lake Storage Path client.
func New() Client {
	return NewWithEnvironment(azure.PublicCloud)
}

// NewWithEnvironment creates an instance of the Data Lake Storage Path client.
func NewWithEnvironment(environment azure.Environment) Client {
	return Client{
		Client:  autorest.NewClientWithUserAgent(UserAgent()),
		BaseURI: environment.StorageEndpointSuffix,
	}
}
