package alertsmanagement

import "github.com/Azure/go-autorest/autorest"

type AlertsManagementClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAlertsManagementClientWithBaseURI(endpoint string) AlertsManagementClient {
	return AlertsManagementClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
