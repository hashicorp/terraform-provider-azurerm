package issue

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IssueClient struct {
	Client *resourcemanager.Client
}

func NewIssueClientWithBaseURI(sdkApi sdkEnv.Api) (*IssueClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "issue", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IssueClient: %+v", err)
	}

	return &IssueClient{
		Client: client,
	}, nil
}
