// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/communicationservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/domains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/emailservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/senderusernames"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SenderUsernamesClient *senderusernames.SenderUsernamesClient
	ServiceClient         *communicationservices.CommunicationServicesClient

	EmailServicesClient *emailservices.EmailServicesClient
	DomainClient        *domains.DomainsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	senderUsernamesClient, err := senderusernames.NewSenderUsernamesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Sender Username client: %+v", err)
	}
	o.Configure(senderUsernamesClient.Client, o.Authorizers.ResourceManager)

	servicesClient, err := communicationservices.NewCommunicationServicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Service client: %+v", err)
	}
	o.Configure(servicesClient.Client, o.Authorizers.ResourceManager)

	emailServicesClient, err := emailservices.NewEmailServicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Email Service client: %+v", err)
	}
	o.Configure(emailServicesClient.Client, o.Authorizers.ResourceManager)

	domainsClient, err := domains.NewDomainsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Domais client: %+v", err)
	}
	o.Configure(domainsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		SenderUsernamesClient: senderUsernamesClient,
		ServiceClient:         servicesClient,
		EmailServicesClient:   emailServicesClient,
		DomainClient:          domainsClient,
	}, nil
}
