package schedule

import "github.com/Azure/go-autorest/autorest"

type ScheduleClient struct {
	Client  autorest.Client
	baseUri string
}

func NewScheduleClientWithBaseURI(endpoint string) ScheduleClient {
	return ScheduleClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
