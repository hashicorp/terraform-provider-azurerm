package v2017_08_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/analysisservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AnalysisServices *analysisservices.AnalysisServicesClient
	Servers          *servers.ServersClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	analysisServicesClient, err := analysisservices.NewAnalysisServicesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AnalysisServices client: %+v", err)
	}
	configureFunc(analysisServicesClient.Client)

	serversClient, err := servers.NewServersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Servers client: %+v", err)
	}
	configureFunc(serversClient.Client)

	return &Client{
		AnalysisServices: analysisServicesClient,
		Servers:          serversClient,
	}, nil
}
