package machinelearningcomputes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachineLearningComputesClient struct {
	Client *resourcemanager.Client
}

func NewMachineLearningComputesClientWithBaseURI(sdkApi sdkEnv.Api) (*MachineLearningComputesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "machinelearningcomputes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MachineLearningComputesClient: %+v", err)
	}

	return &MachineLearningComputesClient{
		Client: client,
	}, nil
}
