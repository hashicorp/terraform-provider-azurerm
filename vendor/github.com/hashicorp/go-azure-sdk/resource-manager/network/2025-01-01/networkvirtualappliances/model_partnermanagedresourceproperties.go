package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerManagedResourceProperties struct {
	Id                     *string `json:"id,omitempty"`
	InternalLoadBalancerId *string `json:"internalLoadBalancerId,omitempty"`
	StandardLoadBalancerId *string `json:"standardLoadBalancerId,omitempty"`
}
