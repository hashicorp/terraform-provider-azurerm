package resolveprivatelinkserviceid

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResolvePrivateLinkServiceIdClient struct {
	Client *resourcemanager.Client
}

func NewResolvePrivateLinkServiceIdClientWithBaseURI(sdkApi sdkEnv.Api) (*ResolvePrivateLinkServiceIdClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "resolveprivatelinkserviceid", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ResolvePrivateLinkServiceIdClient: %+v", err)
	}

	return &ResolvePrivateLinkServiceIdClient{
		Client: client,
	}, nil
}
