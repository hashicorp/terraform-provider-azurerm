package operationstatus

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationStatusClient struct {
	Client  autorest.Client
	baseUri string
}

func NewOperationStatusClientWithBaseURI(endpoint string) OperationStatusClient {
	return OperationStatusClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
