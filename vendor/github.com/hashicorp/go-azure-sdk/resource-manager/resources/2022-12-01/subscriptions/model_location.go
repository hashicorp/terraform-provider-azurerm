package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Location struct {
	AvailabilityZoneMappings *[]AvailabilityZoneMappings `json:"availabilityZoneMappings,omitempty"`
	DisplayName              *string                     `json:"displayName,omitempty"`
	Id                       *string                     `json:"id,omitempty"`
	Metadata                 *LocationMetadata           `json:"metadata,omitempty"`
	Name                     *string                     `json:"name,omitempty"`
	RegionalDisplayName      *string                     `json:"regionalDisplayName,omitempty"`
	SubscriptionId           *string                     `json:"subscriptionId,omitempty"`
	Type                     *LocationType               `json:"type,omitempty"`
}
