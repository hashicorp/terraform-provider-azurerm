package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Location struct {
	DocumentEndpoint  *string `json:"documentEndpoint,omitempty"`
	FailoverPriority  *int64  `json:"failoverPriority,omitempty"`
	Id                *string `json:"id,omitempty"`
	IsZoneRedundant   *bool   `json:"isZoneRedundant,omitempty"`
	LocationName      *string `json:"locationName,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}
