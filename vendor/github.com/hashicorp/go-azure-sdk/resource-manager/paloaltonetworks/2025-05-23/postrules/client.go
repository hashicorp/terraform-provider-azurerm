package postrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PostRulesClient struct {
	Client *resourcemanager.Client
}

func NewPostRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*PostRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "postrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PostRulesClient: %+v", err)
	}

	return &PostRulesClient{
		Client: client,
	}, nil
}
