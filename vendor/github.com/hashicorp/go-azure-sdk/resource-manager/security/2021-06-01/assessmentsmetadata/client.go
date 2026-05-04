package assessmentsmetadata

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssessmentsMetadataClient struct {
	Client *resourcemanager.Client
}

func NewAssessmentsMetadataClientWithBaseURI(sdkApi sdkEnv.Api) (*AssessmentsMetadataClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "assessmentsmetadata", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AssessmentsMetadataClient: %+v", err)
	}

	return &AssessmentsMetadataClient{
		Client: client,
	}, nil
}
