package buckets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BucketsClient struct {
	Client *resourcemanager.Client
}

func NewBucketsClientWithBaseURI(sdkApi sdkEnv.Api) (*BucketsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "buckets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BucketsClient: %+v", err)
	}

	return &BucketsClient{
		Client: client,
	}, nil
}
