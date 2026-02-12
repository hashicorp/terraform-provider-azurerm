package afdcustomdomains

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDCustomDomainsClient struct {
	Client *resourcemanager.Client
}

func NewAFDCustomDomainsClientWithBaseURI(sdkApi sdkEnv.Api) (*AFDCustomDomainsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "afdcustomdomains", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AFDCustomDomainsClient: %+v", err)
	}

	return &AFDCustomDomainsClient{
		Client: client,
	}, nil
}
