package firstpartyservicetags

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirstPartyServiceTagsClient struct {
	Client *resourcemanager.Client
}

func NewFirstPartyServiceTagsClientWithBaseURI(sdkApi sdkEnv.Api) (*FirstPartyServiceTagsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "firstpartyservicetags", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FirstPartyServiceTagsClient: %+v", err)
	}

	return &FirstPartyServiceTagsClient{
		Client: client,
	}, nil
}
