package groupquotalimits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GroupQuotaDetails struct {
	AllocatedToSubscriptions *AllocatedQuotaToSubscriptionList `json:"allocatedToSubscriptions,omitempty"`
	AvailableLimit           *int64                            `json:"availableLimit,omitempty"`
	Comment                  *string                           `json:"comment,omitempty"`
	Limit                    *int64                            `json:"limit,omitempty"`
	Name                     *GroupQuotaDetailsName            `json:"name,omitempty"`
	ResourceName             *string                           `json:"resourceName,omitempty"`
	Unit                     *string                           `json:"unit,omitempty"`
}
