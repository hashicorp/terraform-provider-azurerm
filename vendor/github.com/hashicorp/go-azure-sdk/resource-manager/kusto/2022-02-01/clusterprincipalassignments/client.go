package clusterprincipalassignments

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPrincipalAssignmentsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewClusterPrincipalAssignmentsClientWithBaseURI(endpoint string) ClusterPrincipalAssignmentsClient {
	return ClusterPrincipalAssignmentsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
