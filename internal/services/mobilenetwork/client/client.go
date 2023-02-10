package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/site"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/slice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MobileNetworkClient *mobilenetwork.MobileNetworkClient
	SliceClient         *slice.SliceClient
	SiteClient          *site.SiteClient
}

func NewClient(o *common.ClientOptions) *Client {
	mobileNetworkClient := mobilenetwork.NewMobileNetworkClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&mobileNetworkClient.Client, o.ResourceManagerAuthorizer)

	siteClient := site.NewSiteClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&siteClient.Client, o.ResourceManagerAuthorizer)

	sliceClient := slice.NewSliceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&sliceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MobileNetworkClient: &mobileNetworkClient,
		SiteClient:          &siteClient,
		SliceClient:         &sliceClient,
	}
}
