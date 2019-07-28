package iothub

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2018-12-01-preview/devices"
	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2018-01-22/iothub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ResourceClient       devices.IotHubResourceClient
	DPSResourceClient    iothub.IotDpsResourceClient
	DPSCertificateClient iothub.DpsCertificateClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ResourceClient = devices.NewIotHubResourceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ResourceClient.Client, o.ResourceManagerAuthorizer)

	c.DPSResourceClient = iothub.NewIotDpsResourceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DPSResourceClient.Client, o.ResourceManagerAuthorizer)

	c.DPSCertificateClient = iothub.NewDpsCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DPSCertificateClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
