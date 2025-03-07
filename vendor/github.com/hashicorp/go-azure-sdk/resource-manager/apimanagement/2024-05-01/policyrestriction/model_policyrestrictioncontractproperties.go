package policyrestriction

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyRestrictionContractProperties struct {
	RequireBase *PolicyRestrictionRequireBase `json:"requireBase,omitempty"`
	Scope       *string                       `json:"scope,omitempty"`
}
