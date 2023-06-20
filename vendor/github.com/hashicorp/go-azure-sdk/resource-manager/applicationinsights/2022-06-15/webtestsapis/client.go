package webtestsapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestsAPIsClient struct {
	Client *resourcemanager.Client
}

func NewWebTestsAPIsClientWithBaseURI(api environments.Api) (*WebTestsAPIsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "webtestsapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WebTestsAPIsClient: %+v", err)
	}

	return &WebTestsAPIsClient{
		Client: client,
	}, nil
}
