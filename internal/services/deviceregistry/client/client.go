package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assets"

	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssetClient                *assets.AssetsClient
	AssetEndpointProfileClient *assetendpointprofiles.AssetEndpointProfilesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	assetEndpointProfileClient, err := assetendpointprofiles.NewAssetEndpointProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating AssetEndpointProfiles Client: %+v", err)
	}
	o.Configure(assetEndpointProfileClient.Client, o.Authorizers.ResourceManager)

	assetClient, err := assets.NewAssetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("creating Asset Client: %+v", err)
	}
	o.Configure(assetClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AssetClient:                assetClient,
		AssetEndpointProfileClient: assetEndpointProfileClient,
	}, nil
}
