package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CosmosDbMongoDbApiLinkedServiceTypeProperties struct {
	ConnectionString       string `json:"connectionString"`
	Database               string `json:"database"`
	IsServerVersionAbove32 *bool  `json:"isServerVersionAbove32,omitempty"`
}
