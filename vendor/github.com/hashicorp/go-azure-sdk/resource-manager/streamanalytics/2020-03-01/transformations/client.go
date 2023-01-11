package transformations

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransformationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTransformationsClientWithBaseURI(endpoint string) TransformationsClient {
	return TransformationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
