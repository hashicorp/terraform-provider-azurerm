package webhooks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebHooksClient struct {
	Client *resourcemanager.Client
}

func NewWebHooksClientWithBaseURI(sdkApi sdkEnv.Api) (*WebHooksClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "webhooks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WebHooksClient: %+v", err)
	}

	return &WebHooksClient{
		Client: client,
	}, nil
}
