package tenants

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTenantsClientWithBaseURI(endpoint string) TenantsClient {
	return TenantsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
