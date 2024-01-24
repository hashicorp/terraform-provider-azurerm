package streamingjobs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingJobsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewStreamingJobsClientWithBaseURI(endpoint string) StreamingJobsClient {
	return StreamingJobsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
