package elasticsans

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticSansClient struct {
	Client *resourcemanager.Client
}

func NewElasticSansClientWithBaseURI(sdkApi sdkEnv.Api) (*ElasticSansClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "elasticsans", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ElasticSansClient: %+v", err)
	}

	return &ElasticSansClient{
		Client: client,
	}, nil
}
