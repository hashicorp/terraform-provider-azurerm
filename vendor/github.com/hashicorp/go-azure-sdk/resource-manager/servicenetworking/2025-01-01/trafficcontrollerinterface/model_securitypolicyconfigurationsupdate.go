package trafficcontrollerinterface

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityPolicyConfigurationsUpdate struct {
	WafSecurityPolicy *WafSecurityPolicyUpdate `json:"wafSecurityPolicy,omitempty"`
}
