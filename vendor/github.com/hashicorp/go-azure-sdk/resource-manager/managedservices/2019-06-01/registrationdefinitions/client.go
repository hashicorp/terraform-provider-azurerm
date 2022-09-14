package registrationdefinitions

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistrationDefinitionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRegistrationDefinitionsClientWithBaseURI(endpoint string) RegistrationDefinitionsClient {
	return RegistrationDefinitionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
