package dicomservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DicomServicesClient struct {
	Client *resourcemanager.Client
}

func NewDicomServicesClientWithBaseURI(sdkApi sdkEnv.Api) (*DicomServicesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "dicomservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DicomServicesClient: %+v", err)
	}

	return &DicomServicesClient{
		Client: client,
	}, nil
}
