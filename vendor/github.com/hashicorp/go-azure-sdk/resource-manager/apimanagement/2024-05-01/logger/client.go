package logger

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoggerClient struct {
	Client *resourcemanager.Client
}

func NewLoggerClientWithBaseURI(sdkApi sdkEnv.Api) (*LoggerClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "logger", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LoggerClient: %+v", err)
	}

	return &LoggerClient{
		Client: client,
	}, nil
}
