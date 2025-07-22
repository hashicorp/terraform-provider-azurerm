package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssetsClient                *assets.AssetsClient
	AssetEndpointProfilesClient *assetendpointprofiles.AssetEndpointProfilesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	assetsClient, err := assets.NewAssetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Assets Client: %+v", err)
	}
	o.Configure(assetsClient.Client, o.Authorizers.ResourceManager)

	assetEndpointProfilesClient, err := assetendpointprofiles.NewAssetEndpointProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Asset Endpoint Profiles Client: %+v", err)
	}
	o.Configure(assetEndpointProfilesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AssetsClient:                assetsClient,
		AssetEndpointProfilesClient: assetEndpointProfilesClient,
	}, nil
}
