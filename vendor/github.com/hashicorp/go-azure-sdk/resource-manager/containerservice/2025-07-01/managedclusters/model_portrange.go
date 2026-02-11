package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortRange struct {
	PortEnd   *int64    `json:"portEnd,omitempty"`
	PortStart *int64    `json:"portStart,omitempty"`
	Protocol  *Protocol `json:"protocol,omitempty"`
}
