package extensionoperationstatus

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionOperationStatusClient struct {
	Client  autorest.Client
	baseUri string
}

func NewExtensionOperationStatusClientWithBaseURI(endpoint string) ExtensionOperationStatusClient {
	return ExtensionOperationStatusClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
