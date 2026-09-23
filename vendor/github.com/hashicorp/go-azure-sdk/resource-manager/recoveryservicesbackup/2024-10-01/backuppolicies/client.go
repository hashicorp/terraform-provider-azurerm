package backuppolicies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewBackupPoliciesClientWithBaseURI(endpoint string) BackupPoliciesClient {
	return BackupPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
