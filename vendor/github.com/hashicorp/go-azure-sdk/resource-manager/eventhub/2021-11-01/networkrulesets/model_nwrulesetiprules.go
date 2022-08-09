package networkrulesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NWRuleSetIPRules struct {
	Action *NetworkRuleIPAction `json:"action,omitempty"`
	IPMask *string              `json:"ipMask,omitempty"`
}
