package loadbalancers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosSettings struct {
	DdosProtectionPlan *SubResource                `json:"ddosProtectionPlan,omitempty"`
	ProtectionMode     *DdosSettingsProtectionMode `json:"protectionMode,omitempty"`
}
