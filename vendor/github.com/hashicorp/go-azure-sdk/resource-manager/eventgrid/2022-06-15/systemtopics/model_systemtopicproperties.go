package systemtopics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemTopicProperties struct {
	MetricResourceId  *string                    `json:"metricResourceId,omitempty"`
	ProvisioningState *ResourceProvisioningState `json:"provisioningState,omitempty"`
	Source            *string                    `json:"source,omitempty"`
	TopicType         *string                    `json:"topicType,omitempty"`
}
