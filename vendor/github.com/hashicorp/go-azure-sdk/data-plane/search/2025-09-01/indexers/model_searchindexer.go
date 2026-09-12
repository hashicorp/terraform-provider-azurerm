package indexers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexer struct {
	DataSourceName      string                       `json:"dataSourceName"`
	Description         *string                      `json:"description,omitempty"`
	Disabled            *bool                        `json:"disabled,omitempty"`
	EncryptionKey       *SearchResourceEncryptionKey `json:"encryptionKey,omitempty"`
	FieldMappings       *[]FieldMapping              `json:"fieldMappings,omitempty"`
	Name                string                       `json:"name"`
	OdataEtag           *string                      `json:"@odata.etag,omitempty"`
	OutputFieldMappings *[]FieldMapping              `json:"outputFieldMappings,omitempty"`
	Parameters          *IndexingParameters          `json:"parameters,omitempty"`
	Schedule            *IndexingSchedule            `json:"schedule,omitempty"`
	SkillsetName        *string                      `json:"skillsetName,omitempty"`
	TargetIndexName     string                       `json:"targetIndexName"`
}
