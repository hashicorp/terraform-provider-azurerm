package v2024_09_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapapplicationserverinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapcentralserverinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapdatabaseinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	SAPApplicationServerInstances *sapapplicationserverinstances.SAPApplicationServerInstancesClient
	SAPCentralServerInstances     *sapcentralserverinstances.SAPCentralServerInstancesClient
	SAPDatabaseInstances          *sapdatabaseinstances.SAPDatabaseInstancesClient
	SAPVirtualInstances           *sapvirtualinstances.SAPVirtualInstancesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	sAPApplicationServerInstancesClient, err := sapapplicationserverinstances.NewSAPApplicationServerInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPApplicationServerInstances client: %+v", err)
	}
	configureFunc(sAPApplicationServerInstancesClient.Client)

	sAPCentralServerInstancesClient, err := sapcentralserverinstances.NewSAPCentralServerInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPCentralServerInstances client: %+v", err)
	}
	configureFunc(sAPCentralServerInstancesClient.Client)

	sAPDatabaseInstancesClient, err := sapdatabaseinstances.NewSAPDatabaseInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPDatabaseInstances client: %+v", err)
	}
	configureFunc(sAPDatabaseInstancesClient.Client)

	sAPVirtualInstancesClient, err := sapvirtualinstances.NewSAPVirtualInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPVirtualInstances client: %+v", err)
	}
	configureFunc(sAPVirtualInstancesClient.Client)

	return &Client{
		SAPApplicationServerInstances: sAPApplicationServerInstancesClient,
		SAPCentralServerInstances:     sAPCentralServerInstancesClient,
		SAPDatabaseInstances:          sAPDatabaseInstancesClient,
		SAPVirtualInstances:           sAPVirtualInstancesClient,
	}, nil
}
