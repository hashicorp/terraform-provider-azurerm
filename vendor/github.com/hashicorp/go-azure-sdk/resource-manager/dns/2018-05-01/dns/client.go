package dns

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsClient struct {
	Client *resourcemanager.Client
}

func NewDnsClientWithBaseURI(api environments.Api) (*DnsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "dns", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DnsClient: %+v", err)
	}

	return &DnsClient{
		Client: client,
	}, nil
}
