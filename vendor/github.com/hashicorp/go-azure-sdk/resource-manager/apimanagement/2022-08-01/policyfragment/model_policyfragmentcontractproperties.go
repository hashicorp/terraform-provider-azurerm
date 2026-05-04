package policyfragment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyFragmentContractProperties struct {
	Description *string                      `json:"description,omitempty"`
	Format      *PolicyFragmentContentFormat `json:"format,omitempty"`
	Value       string                       `json:"value"`
}
