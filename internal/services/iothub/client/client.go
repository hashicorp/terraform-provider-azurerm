package client

import (
	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2021-07-02/devices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/dpscertificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ResourceClient          *devices.IotHubResourceClient
	IotHubCertificateClient *devices.CertificatesClient
	DeviceUpdatesClient     *deviceupdates.DeviceupdatesClient
	DPSResourceClient       *iotdpsresource.IotDpsResourceClient
	DPSCertificateClient    *dpscertificate.DpsCertificateClient
}

func NewClient(o *common.ClientOptions) *Client {
	ResourceClient := devices.NewIotHubResourceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ResourceClient.Client, o.ResourceManagerAuthorizer)

	IotHubCertificateClient := devices.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&IotHubCertificateClient.Client, o.ResourceManagerAuthorizer)

	DeviceUpdatesClient := deviceupdates.NewDeviceupdatesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DeviceUpdatesClient.Client, o.ResourceManagerAuthorizer)

	DPSResourceClient := iotdpsresource.NewIotDpsResourceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DPSResourceClient.Client, o.ResourceManagerAuthorizer)

	DPSCertificateClient := dpscertificate.NewDpsCertificateClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DPSCertificateClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ResourceClient:          &ResourceClient,
		IotHubCertificateClient: &IotHubCertificateClient,
		DeviceUpdatesClient:     &DeviceUpdatesClient,
		DPSResourceClient:       &DPSResourceClient,
		DPSCertificateClient:    &DPSCertificateClient,
	}
}
