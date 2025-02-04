package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbV2LinkedServiceTypeProperties struct {
	ConnectionString string `json:"connectionString"`
	Database         string `json:"database"`
}
