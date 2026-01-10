package report

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReportClient struct {
	Client *resourcemanager.Client
}

func NewReportClientWithBaseURI(sdkApi sdkEnv.Api) (*ReportClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "report", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ReportClient: %+v", err)
	}

	return &ReportClient{
		Client: client,
	}, nil
}
