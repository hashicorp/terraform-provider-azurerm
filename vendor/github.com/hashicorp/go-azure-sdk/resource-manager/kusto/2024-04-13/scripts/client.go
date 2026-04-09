package scripts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptsClient struct {
	Client *resourcemanager.Client
}

func NewScriptsClientWithBaseURI(sdkApi sdkEnv.Api) (*ScriptsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "scripts", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ScriptsClient: %+v", err)
	}

	return &ScriptsClient{
		Client: client,
	}, nil
}
