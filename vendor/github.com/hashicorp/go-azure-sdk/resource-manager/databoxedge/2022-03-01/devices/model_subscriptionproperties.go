package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionProperties struct {
	LocationPlacementId *string                           `json:"locationPlacementId,omitempty"`
	QuotaId             *string                           `json:"quotaId,omitempty"`
	RegisteredFeatures  *[]SubscriptionRegisteredFeatures `json:"registeredFeatures,omitempty"`
	SerializedDetails   *string                           `json:"serializedDetails,omitempty"`
	TenantId            *string                           `json:"tenantId,omitempty"`
}
