package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/imageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/pools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/resourcedetails"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/sku"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/subscriptionusages"

	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ImageVersionsClient      *imageversions.ImageVersionsClient
	PoolsClient              *pools.PoolsClient
	ResourceDetailsClient    *resourcedetails.ResourceDetailsClient
	SkuClient                *sku.SkuClient
	SubscriptionUsagesClient *subscriptionusages.SubscriptionUsagesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	ImageVersionsClient, err := imageversions.NewImageVersionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Image Versions Client: %+v", err)
	}
	o.Configure(ImageVersionsClient.Client, o.Authorizers.ResourceManager)

	PoolsClient, err := pools.NewPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Pools Client: %+v", err)
	}
	o.Configure(PoolsClient.Client, o.Authorizers.ResourceManager)

	ResourceDetailsClient, err := resourcedetails.NewResourceDetailsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Resource Details Client: %+v", err)
	}
	o.Configure(ResourceDetailsClient.Client, o.Authorizers.ResourceManager)

	SkuClient, err := sku.NewSkuClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Sku Client: %+v", err)
	}
	o.Configure(SkuClient.Client, o.Authorizers.ResourceManager)

	SubscriptionUsagesClient, err := subscriptionusages.NewSubscriptionUsagesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Subscription Usages Client: %+v", err)
	}
	o.Configure(SubscriptionUsagesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ImageVersionsClient:      ImageVersionsClient,
		PoolsClient:              PoolsClient,
		ResourceDetailsClient:    ResourceDetailsClient,
		SkuClient:                SkuClient,
		SubscriptionUsagesClient: SubscriptionUsagesClient,
	}, nil
}
