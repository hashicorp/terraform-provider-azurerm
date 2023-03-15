package patchschedules

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchSchedulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPatchSchedulesClientWithBaseURI(endpoint string) PatchSchedulesClient {
	return PatchSchedulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
