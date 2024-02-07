package v2022_07_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-07-01/encodings"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Encodings *encodings.EncodingsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	encodingsClient, err := encodings.NewEncodingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Encodings client: %+v", err)
	}
	configureFunc(encodingsClient.Client)

	return &Client{
		Encodings: encodingsClient,
	}, nil
}
