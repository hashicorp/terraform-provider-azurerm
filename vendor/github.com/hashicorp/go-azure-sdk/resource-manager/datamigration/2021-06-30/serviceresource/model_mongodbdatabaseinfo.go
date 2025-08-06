package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbDatabaseInfo struct {
	AverageDocumentSize int64                   `json:"averageDocumentSize"`
	Collections         []MongoDbCollectionInfo `json:"collections"`
	DataSize            int64                   `json:"dataSize"`
	DocumentCount       int64                   `json:"documentCount"`
	Name                string                  `json:"name"`
	QualifiedName       string                  `json:"qualifiedName"`
	SupportsSharding    bool                    `json:"supportsSharding"`
}
