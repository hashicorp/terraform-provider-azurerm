package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsConfiguration struct {
	NameServers   []string `json:"nameServers"`
	Options       *string  `json:"options,omitempty"`
	SearchDomains *string  `json:"searchDomains,omitempty"`
}
