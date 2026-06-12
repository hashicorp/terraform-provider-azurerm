package autonomousdatabases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionURLType struct {
	ApexURL                    *string `json:"apexUrl,omitempty"`
	DatabaseTransformsURL      *string `json:"databaseTransformsUrl,omitempty"`
	GraphStudioURL             *string `json:"graphStudioUrl,omitempty"`
	MachineLearningNotebookURL *string `json:"machineLearningNotebookUrl,omitempty"`
	MongoDbURL                 *string `json:"mongoDbUrl,omitempty"`
	OrdsURL                    *string `json:"ordsUrl,omitempty"`
	SqlDevWebURL               *string `json:"sqlDevWebUrl,omitempty"`
}
