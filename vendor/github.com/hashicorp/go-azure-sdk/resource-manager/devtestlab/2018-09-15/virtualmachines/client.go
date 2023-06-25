package virtualmachines

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachinesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVirtualMachinesClientWithBaseURI(endpoint string) VirtualMachinesClient {
	return VirtualMachinesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
