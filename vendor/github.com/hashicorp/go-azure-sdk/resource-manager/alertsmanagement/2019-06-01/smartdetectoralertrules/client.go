package smartdetectoralertrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SmartDetectorAlertRulesClient struct {
	Client *resourcemanager.Client
}

func NewSmartDetectorAlertRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*SmartDetectorAlertRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "smartdetectoralertrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SmartDetectorAlertRulesClient: %+v", err)
	}

	return &SmartDetectorAlertRulesClient{
		Client: client,
	}, nil
}
