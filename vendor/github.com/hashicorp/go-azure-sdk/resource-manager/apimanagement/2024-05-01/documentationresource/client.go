package documentationresource

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DocumentationResourceClient struct {
	Client *resourcemanager.Client
}

func NewDocumentationResourceClientWithBaseURI(sdkApi sdkEnv.Api) (*DocumentationResourceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "documentationresource", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DocumentationResourceClient: %+v", err)
	}

	return &DocumentationResourceClient{
		Client: client,
	}, nil
}
