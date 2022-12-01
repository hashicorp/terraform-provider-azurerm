package simpolicy

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SIMPolicyClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSIMPolicyClientWithBaseURI(endpoint string) SIMPolicyClient {
	return SIMPolicyClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
