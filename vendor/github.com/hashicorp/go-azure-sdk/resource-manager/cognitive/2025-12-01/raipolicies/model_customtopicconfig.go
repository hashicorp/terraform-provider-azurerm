package raipolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomTopicConfig struct {
	Blocking  *bool                   `json:"blocking,omitempty"`
	Source    *RaiPolicyContentSource `json:"source,omitempty"`
	TopicName *string                 `json:"topicName,omitempty"`
}
