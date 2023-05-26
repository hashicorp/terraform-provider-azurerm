package actiongroupsapis

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsAPIsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewActionGroupsAPIsClientWithBaseURI(endpoint string) ActionGroupsAPIsClient {
	return ActionGroupsAPIsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
