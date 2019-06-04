package iothub

import "github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2018-12-01-preview/devices"

type Client struct {
	ResourceClient devices.IotHubResourceClient
}
