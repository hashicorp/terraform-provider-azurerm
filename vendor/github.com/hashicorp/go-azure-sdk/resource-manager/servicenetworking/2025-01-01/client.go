package v2025_01_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/associationsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/frontendsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/securitypoliciesinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/trafficcontrollerinterface"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AssociationsInterface      *associationsinterface.AssociationsInterfaceClient
	FrontendsInterface         *frontendsinterface.FrontendsInterfaceClient
	SecurityPoliciesInterface  *securitypoliciesinterface.SecurityPoliciesInterfaceClient
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

	securityPoliciesInterfaceClient, err := securitypoliciesinterface.NewSecurityPoliciesInterfaceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SecurityPoliciesInterface client: %+v", err)
	}
	configureFunc(securityPoliciesInterfaceClient.Client)

	trafficControllerInterfaceClient, err := trafficcontrollerinterface.NewTrafficControllerInterfaceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TrafficControllerInterface client: %+v", err)
	}
	configureFunc(trafficControllerInterfaceClient.Client)

	return &Client{
		AssociationsInterface:      associationsInterfaceClient,
		FrontendsInterface:         frontendsInterfaceClient,
		SecurityPoliciesInterface:  securityPoliciesInterfaceClient,
		TrafficControllerInterface: trafficControllerInterfaceClient,
	}, nil
}
