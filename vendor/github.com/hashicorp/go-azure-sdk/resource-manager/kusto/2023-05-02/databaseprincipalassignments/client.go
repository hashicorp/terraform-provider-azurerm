package databaseprincipalassignments

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabasePrincipalAssignmentsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDatabasePrincipalAssignmentsClientWithBaseURI(endpoint string) DatabasePrincipalAssignmentsClient {
	return DatabasePrincipalAssignmentsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
