package datadogsinglesignonresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogSingleSignOnResourcesClient struct {
	Client *resourcemanager.Client
}

func NewDatadogSingleSignOnResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*DatadogSingleSignOnResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "datadogsinglesignonresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DatadogSingleSignOnResourcesClient: %+v", err)
	}

	return &DatadogSingleSignOnResourcesClient{
		Client: client,
	}, nil
}
