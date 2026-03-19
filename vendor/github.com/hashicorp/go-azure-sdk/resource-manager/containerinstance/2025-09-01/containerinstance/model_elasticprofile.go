package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticProfile struct {
	ContainerGroupNamingPolicy *ElasticProfileContainerGroupNamingPolicy `json:"containerGroupNamingPolicy,omitempty"`
	DesiredCount               *int64                                    `json:"desiredCount,omitempty"`
	MaintainDesiredCount       *bool                                     `json:"maintainDesiredCount,omitempty"`
}
