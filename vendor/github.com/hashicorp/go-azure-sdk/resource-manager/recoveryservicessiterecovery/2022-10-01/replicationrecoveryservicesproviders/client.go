package replicationrecoveryservicesproviders

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationRecoveryServicesProvidersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewReplicationRecoveryServicesProvidersClientWithBaseURI(endpoint string) ReplicationRecoveryServicesProvidersClient {
	return ReplicationRecoveryServicesProvidersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
