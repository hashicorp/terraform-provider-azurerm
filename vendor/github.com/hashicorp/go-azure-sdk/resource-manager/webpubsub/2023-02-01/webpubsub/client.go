package webpubsub

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebPubSubClient struct {
	Client *resourcemanager.Client
}

func NewWebPubSubClientWithBaseURI(sdkApi sdkEnv.Api) (*WebPubSubClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "webpubsub", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WebPubSubClient: %+v", err)
	}

	return &WebPubSubClient{
		Client: client,
	}, nil
}
