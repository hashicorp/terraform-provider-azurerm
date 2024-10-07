// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

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

func NewClient(o *common.ClientOptions) (*Client, error) {
	spacecraftClient, err := spacecraft.NewSpacecraftClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Orbital Spacecraft client: %+v", err)
	}
	o.Configure(spacecraftClient.Client, o.Authorizers.ResourceManager)

	contactProfileClient, err := contactprofile.NewContactProfileClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Orbital Contact Profile client: %+v", err)
	}
	o.Configure(contactProfileClient.Client, o.Authorizers.ResourceManager)

	contactClient, err := contact.NewContactClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Orbital Contact client: %+v", err)
	}
	o.Configure(contactClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		SpacecraftClient:     spacecraftClient,
		ContactProfileClient: contactProfileClient,
		ContactClient:        contactClient,
	}, nil
}
