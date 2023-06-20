package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contact"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contactprofile"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/groundstation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/spacecraft"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SpacecraftClient     *spacecraft.SpacecraftClient
	ContactProfileClient *contactprofile.ContactProfileClient
	ContactClient        *contact.ContactClient
	GroundStationClient  *groundstation.GroundStationClient
}

func NewClient(o *common.ClientOptions) *Client {
	spacecraftClient := spacecraft.NewSpacecraftClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&spacecraftClient.Client, o.ResourceManagerAuthorizer)

	contactProfileClient := contactprofile.NewContactProfileClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&contactProfileClient.Client, o.ResourceManagerAuthorizer)

	contactClient := contact.NewContactClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&contactClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SpacecraftClient:     &spacecraftClient,
		ContactProfileClient: &contactProfileClient,
		ContactClient:        &contactClient,
	}
}
