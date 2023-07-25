package scripts

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewScriptsClientWithBaseURI(endpoint string) ScriptsClient {
	return ScriptsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
