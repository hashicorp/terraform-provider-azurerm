package postrulesresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PostRulesResourcesClient struct {
	Client *resourcemanager.Client
}

func NewPostRulesResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*PostRulesResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "postrulesresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PostRulesResourcesClient: %+v", err)
	}

	return &PostRulesResourcesClient{
		Client: client,
	}, nil
}
