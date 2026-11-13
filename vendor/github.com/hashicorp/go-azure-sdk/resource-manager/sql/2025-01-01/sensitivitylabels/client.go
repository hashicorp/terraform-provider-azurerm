package sensitivitylabels

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SensitivityLabelsClient struct {
	Client *resourcemanager.Client
}

func NewSensitivityLabelsClientWithBaseURI(sdkApi sdkEnv.Api) (*SensitivityLabelsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "sensitivitylabels", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SensitivityLabelsClient: %+v", err)
	}

	return &SensitivityLabelsClient{
		Client: client,
	}, nil
}
