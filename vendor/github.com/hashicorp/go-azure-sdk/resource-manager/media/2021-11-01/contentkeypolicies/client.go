package contentkeypolicies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewContentKeyPoliciesClientWithBaseURI(endpoint string) ContentKeyPoliciesClient {
	return ContentKeyPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
