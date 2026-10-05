package templatespecversions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TemplateSpecVersionsClient struct {
	Client *resourcemanager.Client
}

func NewTemplateSpecVersionsClientWithBaseURI(sdkApi sdkEnv.Api) (*TemplateSpecVersionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "templatespecversions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TemplateSpecVersionsClient: %+v", err)
	}

	return &TemplateSpecVersionsClient{
		Client: client,
	}, nil
}
