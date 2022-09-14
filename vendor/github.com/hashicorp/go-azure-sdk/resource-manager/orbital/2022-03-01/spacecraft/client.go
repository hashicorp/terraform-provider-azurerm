package spacecraft

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpacecraftClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSpacecraftClientWithBaseURI(endpoint string) SpacecraftClient {
	return SpacecraftClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
