package apiissue

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiIssueClient struct {
	Client *resourcemanager.Client
}

func NewApiIssueClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiIssueClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "apiissue", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiIssueClient: %+v", err)
	}

	return &ApiIssueClient{
		Client: client,
	}, nil
}
