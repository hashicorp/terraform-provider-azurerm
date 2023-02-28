package refreshsetpasswordlink

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RefreshSetPasswordLinkClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRefreshSetPasswordLinkClientWithBaseURI(endpoint string) RefreshSetPasswordLinkClient {
	return RefreshSetPasswordLinkClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
