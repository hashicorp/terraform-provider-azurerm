package registrationassignments

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistrationAssignmentsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRegistrationAssignmentsClientWithBaseURI(endpoint string) RegistrationAssignmentsClient {
	return RegistrationAssignmentsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
