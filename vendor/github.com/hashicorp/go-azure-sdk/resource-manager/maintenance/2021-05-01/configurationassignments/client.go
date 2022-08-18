package configurationassignments

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationAssignmentsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConfigurationAssignmentsClientWithBaseURI(endpoint string) ConfigurationAssignmentsClient {
	return ConfigurationAssignmentsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
