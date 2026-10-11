package blobservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobServicesClient struct {
	Client *resourcemanager.Client
}

func NewBlobServicesClientWithBaseURI(sdkApi sdkEnv.Api) (*BlobServicesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "blobservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BlobServicesClient: %+v", err)
	}

	return &BlobServicesClient{
		Client: client,
	}, nil
}
