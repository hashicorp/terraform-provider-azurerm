package outputs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubOutputDataSourceProperties struct {
	AuthenticationMode     *AuthenticationMode `json:"authenticationMode,omitempty"`
	EventHubName           *string             `json:"eventHubName,omitempty"`
	PartitionCount         *int64              `json:"partitionCount,omitempty"`
	PartitionKey           *string             `json:"partitionKey,omitempty"`
	PropertyColumns        *[]string           `json:"propertyColumns,omitempty"`
	ServiceBusNamespace    *string             `json:"serviceBusNamespace,omitempty"`
	SharedAccessPolicyKey  *string             `json:"sharedAccessPolicyKey,omitempty"`
	SharedAccessPolicyName *string             `json:"sharedAccessPolicyName,omitempty"`
}
