package azurebackupjobs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBackupJobsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAzureBackupJobsClientWithBaseURI(endpoint string) AzureBackupJobsClient {
	return AzureBackupJobsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
