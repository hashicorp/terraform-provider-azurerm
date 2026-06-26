package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrderBy struct {
	Field *string                           `json:"field,omitempty"`
	Order *FirewallPolicyIDPSQuerySortOrder `json:"order,omitempty"`
}
