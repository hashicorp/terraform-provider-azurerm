package networkanchors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsForwardingRule struct {
	DomainNames         string `json:"domainNames"`
	ForwardingIPAddress string `json:"forwardingIpAddress"`
}
