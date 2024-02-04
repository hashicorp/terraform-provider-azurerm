package billingaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndirectRelationshipInfo struct {
	BillingAccountName *string `json:"billingAccountName,omitempty"`
	BillingProfileName *string `json:"billingProfileName,omitempty"`
	DisplayName        *string `json:"displayName,omitempty"`
}
