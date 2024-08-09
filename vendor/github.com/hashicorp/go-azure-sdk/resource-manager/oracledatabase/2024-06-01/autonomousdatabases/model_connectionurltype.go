package autonomousdatabases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionUrlType struct {
	ApexUrl                    *string `json:"apexUrl,omitempty"`
	DatabaseTransformsUrl      *string `json:"databaseTransformsUrl,omitempty"`
	GraphStudioUrl             *string `json:"graphStudioUrl,omitempty"`
	MachineLearningNotebookUrl *string `json:"machineLearningNotebookUrl,omitempty"`
	MongoDbUrl                 *string `json:"mongoDbUrl,omitempty"`
	OrdsUrl                    *string `json:"ordsUrl,omitempty"`
	SqlDevWebUrl               *string `json:"sqlDevWebUrl,omitempty"`
}
