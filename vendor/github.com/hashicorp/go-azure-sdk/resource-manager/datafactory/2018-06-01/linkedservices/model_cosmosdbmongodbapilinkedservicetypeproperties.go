package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CosmosDbMongoDbApiLinkedServiceTypeProperties struct {
	ConnectionString       interface{} `json:"connectionString"`
	Database               interface{} `json:"database"`
	IsServerVersionAbove32 *bool       `json:"isServerVersionAbove32,omitempty"`
}
