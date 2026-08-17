package threatintelligence

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceClient struct {
	Client *resourcemanager.Client
}

func NewThreatIntelligenceClientWithBaseURI(sdkApi sdkEnv.Api) (*ThreatIntelligenceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "threatintelligence", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ThreatIntelligenceClient: %+v", err)
	}

	return &ThreatIntelligenceClient{
		Client: client,
	}, nil
}
