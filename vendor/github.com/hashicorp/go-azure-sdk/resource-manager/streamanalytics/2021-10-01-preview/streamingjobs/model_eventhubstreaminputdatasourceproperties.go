package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubStreamInputDataSourceProperties struct {
	AuthenticationMode     *AuthenticationMode `json:"authenticationMode,omitempty"`
	ConsumerGroupName      *string             `json:"consumerGroupName,omitempty"`
	EventHubName           *string             `json:"eventHubName,omitempty"`
	PartitionCount         *int64              `json:"partitionCount,omitempty"`
	PrefetchCount          *int64              `json:"prefetchCount,omitempty"`
	ServiceBusNamespace    *string             `json:"serviceBusNamespace,omitempty"`
	SharedAccessPolicyKey  *string             `json:"sharedAccessPolicyKey,omitempty"`
	SharedAccessPolicyName *string             `json:"sharedAccessPolicyName,omitempty"`
}
