// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/attacheddatanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/datanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcorecontrolplane"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcoredataplane"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/service"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/sim"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/site"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/slice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MobileNetworkClient          *mobilenetwork.MobileNetworkClient
	ServiceClient                *service.ServiceClient
	SIMGroupClient               *simgroup.SIMGroupClient
	SliceClient                  *slice.SliceClient
	SiteClient                   *site.SiteClient
	DataNetworkClient            *datanetwork.DataNetworkClient
	SIMPolicyClient              *simpolicy.SIMPolicyClient
	PacketCoreControlPlaneClient *packetcorecontrolplane.PacketCoreControlPlaneClient
	PacketCoreDataPlaneClient    *packetcoredataplane.PacketCoreDataPlaneClient
	AttachedDataNetworkClient    *attacheddatanetwork.AttachedDataNetworkClient
	SIMClient                    *sim.SIMClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	mobileNetworkClient, err := mobilenetwork.NewMobileNetworkClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Mobile Network Client: %+v", err)
	}
	o.Configure(mobileNetworkClient.Client, o.Authorizers.ResourceManager)

	dataNetworkClient, err := datanetwork.NewDataNetworkClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Data Network Client: %+v", err)
	}
	o.Configure(dataNetworkClient.Client, o.Authorizers.ResourceManager)

	serviceClient, err := service.NewServiceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Service Client: %+v", err)
	}
	o.Configure(serviceClient.Client, o.Authorizers.ResourceManager)

	simGroupClient, err := simgroup.NewSIMGroupClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SIM Group Client: %+v", err)
	}
	o.Configure(simGroupClient.Client, o.Authorizers.ResourceManager)

	siteClient, err := site.NewSiteClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Site Client: %+v", err)
	}
	o.Configure(siteClient.Client, o.Authorizers.ResourceManager)

	sliceClient, err := slice.NewSliceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Slice Client: %+v", err)
	}
	o.Configure(sliceClient.Client, o.Authorizers.ResourceManager)

	simPolicyClient, err := simpolicy.NewSIMPolicyClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SIM Policy Client: %+v", err)
	}
	o.Configure(simPolicyClient.Client, o.Authorizers.ResourceManager)

	packetCoreControlPlaneClient, err := packetcorecontrolplane.NewPacketCoreControlPlaneClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Packet Core Control Plane Client: %+v", err)
	}
	o.Configure(packetCoreControlPlaneClient.Client, o.Authorizers.ResourceManager)

	packetCoreDataPlaneClient, err := packetcoredataplane.NewPacketCoreDataPlaneClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Packet Core Data Plane Client: %+v", err)
	}
	o.Configure(packetCoreDataPlaneClient.Client, o.Authorizers.ResourceManager)

	attachedDataNetworkClient, err := attacheddatanetwork.NewAttachedDataNetworkClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Attached Data Network Client: %+v", err)
	}
	o.Configure(attachedDataNetworkClient.Client, o.Authorizers.ResourceManager)

	simClient, err := sim.NewSIMClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SIM Client: %+v", err)
	}
	o.Configure(simClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MobileNetworkClient:          mobileNetworkClient,
		DataNetworkClient:            dataNetworkClient,
		ServiceClient:                serviceClient,
		SIMGroupClient:               simGroupClient,
		SiteClient:                   siteClient,
		SliceClient:                  sliceClient,
		SIMPolicyClient:              simPolicyClient,
		PacketCoreControlPlaneClient: packetCoreControlPlaneClient,
		PacketCoreDataPlaneClient:    packetCoreDataPlaneClient,
		AttachedDataNetworkClient:    attachedDataNetworkClient,
		SIMClient:                    simClient,
	}, nil
}
