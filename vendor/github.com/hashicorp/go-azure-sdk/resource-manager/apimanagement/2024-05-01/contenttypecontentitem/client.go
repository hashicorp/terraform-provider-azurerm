package contenttypecontentitem

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentTypeContentItemClient struct {
	Client *resourcemanager.Client
}

func NewContentTypeContentItemClientWithBaseURI(sdkApi sdkEnv.Api) (*ContentTypeContentItemClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "contenttypecontentitem", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ContentTypeContentItemClient: %+v", err)
	}

	return &ContentTypeContentItemClient{
		Client: client,
	}, nil
}
