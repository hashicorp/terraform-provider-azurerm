package topicrecords

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicProperties struct {
	Configs                 *TopicsRelatedLink   `json:"configs,omitempty"`
	InputConfigs            *[]TopicsInputConfig `json:"inputConfigs,omitempty"`
	Kind                    *string              `json:"kind,omitempty"`
	Metadata                *TopicMetadataEntity `json:"metadata,omitempty"`
	Partitions              *TopicsRelatedLink   `json:"partitions,omitempty"`
	PartitionsCount         *string              `json:"partitionsCount,omitempty"`
	PartitionsReassignments *TopicsRelatedLink   `json:"partitionsReassignments,omitempty"`
	ReplicationFactor       *string              `json:"replicationFactor,omitempty"`
	TopicId                 *string              `json:"topicId,omitempty"`
}
