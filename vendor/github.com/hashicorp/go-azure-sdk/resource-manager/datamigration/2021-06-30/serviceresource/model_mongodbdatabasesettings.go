package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbDatabaseSettings struct {
	Collections map[string]MongoDbCollectionSettings `json:"collections"`
	TargetRUs   *int64                               `json:"targetRUs,omitempty"`
}
