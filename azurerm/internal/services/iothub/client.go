package iothub

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2018-12-01-preview/devices"
	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2018-01-22/iothub"
)

type Client struct {
	ResourceClient       devices.IotHubResourceClient
	DPSResourceClient    iothub.IotDpsResourceClient
	DPSCertificateClient iothub.DpsCertificateClient
}
