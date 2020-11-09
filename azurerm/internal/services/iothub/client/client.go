package client

import (
	"github.com/Azure/azure-sdk-for-go/services/digitaltwins/mgmt/2020-10-31/digitaltwins"
	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2020-03-01/devices"
	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2018-01-22/iothub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ResourceClient       *devices.IotHubResourceClient
	DPSResourceClient    *iothub.IotDpsResourceClient
	DPSCertificateClient *iothub.DpsCertificateClient
	DigitalTwinsClient   *digitaltwins.Client
}

func NewClient(o *common.ClientOptions) *Client {
	ResourceClient := devices.NewIotHubResourceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ResourceClient.Client, o.ResourceManagerAuthorizer)

	DPSResourceClient := iothub.NewIotDpsResourceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DPSResourceClient.Client, o.ResourceManagerAuthorizer)

	DPSCertificateClient := iothub.NewDpsCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DPSCertificateClient.Client, o.ResourceManagerAuthorizer)

	digitalTwinsClient := digitaltwins.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&digitalTwinsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ResourceClient:       &ResourceClient,
		DPSResourceClient:    &DPSResourceClient,
		DPSCertificateClient: &DPSCertificateClient,
		DigitalTwinsClient:   &digitalTwinsClient,
	}
}
