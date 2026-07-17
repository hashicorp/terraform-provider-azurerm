package runbookdraft

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunbookDraftClient struct {
	Client *resourcemanager.Client
}

func NewRunbookDraftClientWithBaseURI(sdkApi sdkEnv.Api) (*RunbookDraftClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "runbookdraft", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RunbookDraftClient: %+v", err)
	}

	return &RunbookDraftClient{
		Client: client,
	}, nil
}
