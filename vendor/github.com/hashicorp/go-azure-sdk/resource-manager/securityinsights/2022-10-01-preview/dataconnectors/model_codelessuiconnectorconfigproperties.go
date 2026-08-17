package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CodelessUiConnectorConfigProperties struct {
	Availability          Availability               `json:"availability"`
	ConnectivityCriteria  []ConnectivityCriteria     `json:"connectivityCriteria"`
	CustomImage           *string                    `json:"customImage,omitempty"`
	DataTypes             []LastDataReceivedDataType `json:"dataTypes"`
	DescriptionMarkdown   string                     `json:"descriptionMarkdown"`
	GraphQueries          []GraphQueries             `json:"graphQueries"`
	GraphQueriesTableName string                     `json:"graphQueriesTableName"`
	InstructionSteps      []InstructionSteps         `json:"instructionSteps"`
	Permissions           Permissions                `json:"permissions"`
	Publisher             string                     `json:"publisher"`
	SampleQueries         []SampleQueries            `json:"sampleQueries"`
	Title                 string                     `json:"title"`
}
