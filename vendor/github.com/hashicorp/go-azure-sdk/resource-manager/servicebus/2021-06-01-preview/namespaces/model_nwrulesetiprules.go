package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NWRuleSetIpRules struct {
	Action *NetworkRuleIPAction `json:"action,omitempty"`
	IpMask *string              `json:"ipMask,omitempty"`
}
