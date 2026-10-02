package blobauditing

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobAuditingClient struct {
	Client *resourcemanager.Client
}

func NewBlobAuditingClientWithBaseURI(sdkApi sdkEnv.Api) (*BlobAuditingClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "blobauditing", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BlobAuditingClient: %+v", err)
	}

	return &BlobAuditingClient{
		Client: client,
	}, nil
}
