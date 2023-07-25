package v2023_04_13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Graphservicesprods *graphservicesprods.GraphservicesprodsClient
}

func NewClientWithBaseURI(api environments.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	graphservicesprodsClient, err := graphservicesprods.NewGraphservicesprodsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Graphservicesprods client: %+v", err)
	}
	configureFunc(graphservicesprodsClient.Client)

	return &Client{
		Graphservicesprods: graphservicesprodsClient,
	}, nil
}
