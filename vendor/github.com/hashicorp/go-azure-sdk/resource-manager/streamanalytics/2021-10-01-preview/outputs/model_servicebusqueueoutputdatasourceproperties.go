package outputs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceBusQueueOutputDataSourceProperties struct {
	AuthenticationMode     *AuthenticationMode `json:"authenticationMode,omitempty"`
	PropertyColumns        *[]string           `json:"propertyColumns,omitempty"`
	QueueName              *string             `json:"queueName,omitempty"`
	ServiceBusNamespace    *string             `json:"serviceBusNamespace,omitempty"`
	SharedAccessPolicyKey  *string             `json:"sharedAccessPolicyKey,omitempty"`
	SharedAccessPolicyName *string             `json:"sharedAccessPolicyName,omitempty"`
	SystemPropertyColumns  *interface{}        `json:"systemPropertyColumns,omitempty"`
}
