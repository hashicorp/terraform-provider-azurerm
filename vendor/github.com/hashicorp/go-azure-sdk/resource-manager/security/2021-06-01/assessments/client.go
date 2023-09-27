package assessments

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssessmentsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAssessmentsClientWithBaseURI(endpoint string) AssessmentsClient {
	return AssessmentsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
