// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/dpscertificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

type Client struct {
	ResourceClient          *devices.IotHubResourceClient
	IotHubCertificateClient *devices.CertificatesClient
	DeviceUpdatesClient     *deviceupdates.DeviceupdatesClient
	DPSResourceClient       *iotdpsresource.IotDpsResourceClient
	DPSCertificateClient    *dpscertificate.DpsCertificateClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	ResourceClient := devices.NewIotHubResourceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ResourceClient.Client, o.ResourceManagerAuthorizer)

	IotHubCertificateClient := devices.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&IotHubCertificateClient.Client, o.ResourceManagerAuthorizer)

	DeviceUpdatesClient, err := deviceupdates.NewDeviceupdatesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Device Updates Client: %+v", err)
	}
	o.Configure(DeviceUpdatesClient.Client, o.Authorizers.ResourceManager)

	DPSResourceClient, err := iotdpsresource.NewIotDpsResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DPS Resource Client: %+v", err)
	}
	o.Configure(DPSResourceClient.Client, o.Authorizers.ResourceManager)

	DPSCertificateClient, err := dpscertificate.NewDpsCertificateClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DPS Certificate Client: %+v", err)
	}
	o.Configure(DPSCertificateClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ResourceClient:          &ResourceClient,
		IotHubCertificateClient: &IotHubCertificateClient,
		DeviceUpdatesClient:     DeviceUpdatesClient,
		DPSResourceClient:       DPSResourceClient,
		DPSCertificateClient:    DPSCertificateClient,
	}, nil
}
