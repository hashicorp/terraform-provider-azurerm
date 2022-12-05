package v2021_12_01_preview

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2021-12-01-preview/loadtests"
)

type Client struct {
	LoadTests *loadtests.LoadTestsClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	loadTestsClient := loadtests.NewLoadTestsClientWithBaseURI(endpoint)
	configureAuthFunc(&loadTestsClient.Client)

	return Client{
		LoadTests: &loadTestsClient,
	}
}
