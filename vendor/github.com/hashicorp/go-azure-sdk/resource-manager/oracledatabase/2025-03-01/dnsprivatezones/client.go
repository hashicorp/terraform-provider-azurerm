package dnsprivatezones

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsPrivateZonesClient struct {
	Client *resourcemanager.Client
}

func NewDnsPrivateZonesClientWithBaseURI(sdkApi sdkEnv.Api) (*DnsPrivateZonesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dnsprivatezones", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DnsPrivateZonesClient: %+v", err)
	}

	return &DnsPrivateZonesClient{
		Client: client,
	}, nil
}
