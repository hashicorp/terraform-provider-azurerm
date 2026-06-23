package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DualStackEndpointPreference struct {
	PublishIPv6Endpoint *bool `json:"publishIpv6Endpoint,omitempty"`
}
