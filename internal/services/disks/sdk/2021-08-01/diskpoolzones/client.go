package diskpoolzones

import "github.com/Azure/go-autorest/autorest"

type DiskPoolZonesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDiskPoolZonesClientWithBaseURI(endpoint string) DiskPoolZonesClient {
	return DiskPoolZonesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
