package hybridrunbookworkergroup

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridRunbookWorkerGroupClient struct {
	Client  autorest.Client
	baseUri string
}

func NewHybridRunbookWorkerGroupClientWithBaseURI(endpoint string) HybridRunbookWorkerGroupClient {
	return HybridRunbookWorkerGroupClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
