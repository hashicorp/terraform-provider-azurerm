package applicationinsights

import "github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"

type Client struct {
	APIKeyClient     insights.APIKeysClient
	ComponentsClient insights.ComponentsClient
	WebTestsClient   insights.WebTestsClient
}
