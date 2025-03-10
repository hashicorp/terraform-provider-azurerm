package skus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementSkuRestrictions struct {
	ReasonCode      *ApiManagementSkuRestrictionsReasonCode `json:"reasonCode,omitempty"`
	RestrictionInfo *ApiManagementSkuRestrictionInfo        `json:"restrictionInfo,omitempty"`
	Type            *ApiManagementSkuRestrictionsType       `json:"type,omitempty"`
	Values          *[]string                               `json:"values,omitempty"`
}
