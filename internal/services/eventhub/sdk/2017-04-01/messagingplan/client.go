package messagingplan

import "github.com/Azure/go-autorest/autorest"

type MessagingPlanClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMessagingPlanClientWithBaseURI(endpoint string) MessagingPlanClient {
	return MessagingPlanClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
