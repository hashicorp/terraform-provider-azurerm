package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/attacheddatanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/datanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/packetcorecontrolplane"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/packetcoredataplane"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/service"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/sim"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/simgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/simpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/site"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/slice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MobileNetworkClient          *mobilenetwork.MobileNetworkClient
	SIMClient                    *sim.SIMClient
	DataNetworkClient            *datanetwork.DataNetworkClient
	PacketCoreControlPlaneClient *packetcorecontrolplane.PacketCoreControlPlaneClient
	ServiceClient                *service.ServiceClient
	PacketCoreDataPlaneClient    *packetcoredataplane.PacketCoreDataPlaneClient
	SliceClient                  *slice.SliceClient
	SIMGroupClient               *simgroup.SIMGroupClient
	AttachedDataNetworkClient    *attacheddatanetwork.AttachedDataNetworkClient
	SIMPolicyClient              *simpolicy.SIMPolicyClient
	SiteClient                   *site.SiteClient
}

func NewClient(o *common.ClientOptions) *Client {

	mobileNetworkClient := mobilenetwork.NewMobileNetworkClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&mobileNetworkClient.Client, o.ResourceManagerAuthorizer)

	simClient := sim.NewSIMClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&simClient.Client, o.ResourceManagerAuthorizer)

	dataNetworkClient := datanetwork.NewDataNetworkClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dataNetworkClient.Client, o.ResourceManagerAuthorizer)

	packetCoreControlPlaneClient := packetcorecontrolplane.NewPacketCoreControlPlaneClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&packetCoreControlPlaneClient.Client, o.ResourceManagerAuthorizer)

	serviceClient := service.NewServiceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serviceClient.Client, o.ResourceManagerAuthorizer)

	packetCoreDataPlaneClient := packetcoredataplane.NewPacketCoreDataPlaneClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&packetCoreDataPlaneClient.Client, o.ResourceManagerAuthorizer)

	sliceClient := slice.NewSliceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&sliceClient.Client, o.ResourceManagerAuthorizer)

	sIMGroupClient := simgroup.NewSIMGroupClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&sIMGroupClient.Client, o.ResourceManagerAuthorizer)

	attachedDataNetworkClient := attacheddatanetwork.NewAttachedDataNetworkClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&attachedDataNetworkClient.Client, o.ResourceManagerAuthorizer)

	sIMPolicyClient := simpolicy.NewSIMPolicyClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&sIMPolicyClient.Client, o.ResourceManagerAuthorizer)

	siteClient := site.NewSiteClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&siteClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MobileNetworkClient:          &mobileNetworkClient,
		SIMClient:                    &simClient,
		DataNetworkClient:            &dataNetworkClient,
		PacketCoreControlPlaneClient: &packetCoreControlPlaneClient,
		ServiceClient:                &serviceClient,
		PacketCoreDataPlaneClient:    &packetCoreDataPlaneClient,
		SliceClient:                  &sliceClient,
		SIMGroupClient:               &sIMGroupClient,
		AttachedDataNetworkClient:    &attachedDataNetworkClient,
		SIMPolicyClient:              &sIMPolicyClient,
		SiteClient:                   &siteClient,
	}
}
