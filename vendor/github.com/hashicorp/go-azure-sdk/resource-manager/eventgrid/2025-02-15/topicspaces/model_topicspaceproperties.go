package topicspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicSpaceProperties struct {
	Description       *string                      `json:"description,omitempty"`
	ProvisioningState *TopicSpaceProvisioningState `json:"provisioningState,omitempty"`
	TopicTemplates    *[]string                    `json:"topicTemplates,omitempty"`
}
