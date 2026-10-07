package afddomains

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDDomainsClient struct {
	Client *resourcemanager.Client
}

func NewAFDDomainsClientWithBaseURI(sdkApi sdkEnv.Api) (*AFDDomainsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "afddomains", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AFDDomainsClient: %+v", err)
	}

	return &AFDDomainsClient{
		Client: client,
	}, nil
}
