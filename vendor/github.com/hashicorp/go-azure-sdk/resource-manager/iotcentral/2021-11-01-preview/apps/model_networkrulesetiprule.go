package apps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkRuleSetIPRule struct {
	Action     *IPRuleAction `json:"action,omitempty"`
	FilterName *string       `json:"filterName,omitempty"`
	IPMask     *string       `json:"ipMask,omitempty"`
}
