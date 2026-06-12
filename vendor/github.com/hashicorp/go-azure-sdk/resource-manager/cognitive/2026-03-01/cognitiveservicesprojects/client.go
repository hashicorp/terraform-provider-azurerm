package cognitiveservicesprojects

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CognitiveServicesProjectsClient struct {
	Client *resourcemanager.Client
}

func NewCognitiveServicesProjectsClientWithBaseURI(sdkApi sdkEnv.Api) (*CognitiveServicesProjectsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "cognitiveservicesprojects", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CognitiveServicesProjectsClient: %+v", err)
	}

	return &CognitiveServicesProjectsClient{
		Client: client,
	}, nil
}
