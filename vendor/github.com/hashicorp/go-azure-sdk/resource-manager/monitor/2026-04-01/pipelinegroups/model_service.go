package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Service struct {
	Persistence *PersistenceConfigurations `json:"persistence,omitempty"`
	Pipelines   []Pipeline                 `json:"pipelines"`
}
