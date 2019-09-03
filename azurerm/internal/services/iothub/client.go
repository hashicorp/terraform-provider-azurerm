package iothub

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2018-12-01-preview/devices"
	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2018-01-22/iothub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ResourceClient       *devices.IotHubResourceClient
	DPSResourceClient    *iothub.IotDpsResourceClient
	DPSCertificateClient *iothub.DpsCertificateClient
}

func BuildClient(o *common.ClientOptions) *Client {

	ResourceClient := devices.NewIotHubResourceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ResourceClient.Client, o.ResourceManagerAuthorizer)

	DPSResourceClient := iothub.NewIotDpsResourceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DPSResourceClient.Client, o.ResourceManagerAuthorizer)

	DPSCertificateClient := iothub.NewDpsCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DPSCertificateClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ResourceClient:       &ResourceClient,
		DPSResourceClient:    &DPSResourceClient,
		DPSCertificateClient: &DPSCertificateClient,
	}
}
