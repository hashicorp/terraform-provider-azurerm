package azurefirewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HubPublicIPAddresses struct {
	Addresses *[]AzureFirewallPublicIPAddress `json:"addresses,omitempty"`
	Count     *int64                          `json:"count,omitempty"`
}
