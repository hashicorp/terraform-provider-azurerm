package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsConfig struct {
	Fqdn         *string `json:"fqdn,omitempty"`
	RelativeName *string `json:"relativeName,omitempty"`
	Ttl          *int64  `json:"ttl,omitempty"`
}
