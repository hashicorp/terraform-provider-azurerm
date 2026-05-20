package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbAtlasLinkedServiceTypeProperties struct {
	ConnectionString interface{}  `json:"connectionString"`
	Database         interface{}  `json:"database"`
	DriverVersion    *interface{} `json:"driverVersion,omitempty"`
}
