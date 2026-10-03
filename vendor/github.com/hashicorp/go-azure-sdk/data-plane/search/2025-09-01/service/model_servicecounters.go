package service

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceCounters struct {
	DataSourcesCount ResourceCounter `json:"dataSourcesCount"`
	DocumentCount    ResourceCounter `json:"documentCount"`
	IndexersCount    ResourceCounter `json:"indexersCount"`
	IndexesCount     ResourceCounter `json:"indexesCount"`
	SkillsetCount    ResourceCounter `json:"skillsetCount"`
	StorageSize      ResourceCounter `json:"storageSize"`
	SynonymMaps      ResourceCounter `json:"synonymMaps"`
	VectorIndexSize  ResourceCounter `json:"vectorIndexSize"`
}
