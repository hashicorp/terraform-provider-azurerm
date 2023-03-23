package replicas

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicasClient struct {
	Client  autorest.Client
	baseUri string
}

func NewReplicasClientWithBaseURI(endpoint string) ReplicasClient {
	return ReplicasClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
