package skillsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DocumentIntelligenceLayoutSkillChunkingProperties struct {
	MaximumLength *int64                                       `json:"maximumLength,omitempty"`
	OverlapLength *int64                                       `json:"overlapLength,omitempty"`
	Unit          *DocumentIntelligenceLayoutSkillChunkingUnit `json:"unit,omitempty"`
}
