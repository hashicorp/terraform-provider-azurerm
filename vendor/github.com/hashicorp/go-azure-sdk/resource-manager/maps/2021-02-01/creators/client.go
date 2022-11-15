package creators

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreatorsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCreatorsClientWithBaseURI(endpoint string) CreatorsClient {
	return CreatorsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
