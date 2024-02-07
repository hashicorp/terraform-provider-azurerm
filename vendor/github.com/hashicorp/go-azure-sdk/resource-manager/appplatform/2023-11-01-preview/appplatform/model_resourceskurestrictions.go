package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceSkuRestrictions struct {
	ReasonCode      *ResourceSkuRestrictionsReasonCode `json:"reasonCode,omitempty"`
	RestrictionInfo *ResourceSkuRestrictionInfo        `json:"restrictionInfo,omitempty"`
	Type            *ResourceSkuRestrictionsType       `json:"type,omitempty"`
	Values          *[]string                          `json:"values,omitempty"`
}
