package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceBusTopicOutputDataSourceProperties struct {
	AuthenticationMode     *AuthenticationMode `json:"authenticationMode,omitempty"`
	PropertyColumns        *[]string           `json:"propertyColumns,omitempty"`
	ServiceBusNamespace    *string             `json:"serviceBusNamespace,omitempty"`
	SharedAccessPolicyKey  *string             `json:"sharedAccessPolicyKey,omitempty"`
	SharedAccessPolicyName *string             `json:"sharedAccessPolicyName,omitempty"`
	SystemPropertyColumns  *map[string]string  `json:"systemPropertyColumns,omitempty"`
	TopicName              *string             `json:"topicName,omitempty"`
}
