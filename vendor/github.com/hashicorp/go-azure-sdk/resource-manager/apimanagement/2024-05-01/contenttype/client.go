package contenttype

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentTypeClient struct {
	Client *resourcemanager.Client
}

func NewContentTypeClientWithBaseURI(sdkApi sdkEnv.Api) (*ContentTypeClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "contenttype", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ContentTypeClient: %+v", err)
	}

	return &ContentTypeClient{
		Client: client,
	}, nil
}
