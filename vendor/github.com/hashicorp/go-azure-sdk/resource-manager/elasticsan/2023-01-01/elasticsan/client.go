package elasticsan

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticSanClient struct {
	Client *resourcemanager.Client
}

func NewElasticSanClientWithBaseURI(sdkApi sdkEnv.Api) (*ElasticSanClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "elasticsan", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ElasticSanClient: %+v", err)
	}

	return &ElasticSanClient{
		Client: client,
	}, nil
}
