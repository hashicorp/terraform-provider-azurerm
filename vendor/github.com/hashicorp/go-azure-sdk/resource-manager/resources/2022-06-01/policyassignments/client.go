package policyassignments

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyAssignmentsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPolicyAssignmentsClientWithBaseURI(endpoint string) PolicyAssignmentsClient {
	return PolicyAssignmentsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
