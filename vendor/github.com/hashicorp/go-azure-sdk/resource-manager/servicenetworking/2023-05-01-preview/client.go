package v2023_05_01_preview

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/associationsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/frontendsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/trafficcontrollerinterface"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AssociationsInterface      *associationsinterface.AssociationsInterfaceClient
	FrontendsInterface         *frontendsinterface.FrontendsInterfaceClient
	TrafficControllerInterface *trafficcontrollerinterface.TrafficControllerInterfaceClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	associationsInterfaceClient, err := associationsinterface.NewAssociationsInterfaceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AssociationsInterface client: %+v", err)
	}
	configureFunc(associationsInterfaceClient.Client)

	frontendsInterfaceClient, err := frontendsinterface.NewFrontendsInterfaceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FrontendsInterface client: %+v", err)
	}
	configureFunc(frontendsInterfaceClient.Client)

	trafficControllerInterfaceClient, err := trafficcontrollerinterface.NewTrafficControllerInterfaceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TrafficControllerInterface client: %+v", err)
	}
	configureFunc(trafficControllerInterfaceClient.Client)

	return &Client{
		AssociationsInterface:      associationsInterfaceClient,
		FrontendsInterface:         frontendsInterfaceClient,
		TrafficControllerInterface: trafficControllerInterfaceClient,
	}, nil
}
