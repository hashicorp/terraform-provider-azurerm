package v2020_12_08

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2020-12-08/healthbots"
)

type Client struct {
	Healthbots *healthbots.HealthbotsClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	healthbotsClient := healthbots.NewHealthbotsClientWithBaseURI(endpoint)
	configureAuthFunc(&healthbotsClient.Client)

	return Client{
		Healthbots: &healthbotsClient,
	}
}
