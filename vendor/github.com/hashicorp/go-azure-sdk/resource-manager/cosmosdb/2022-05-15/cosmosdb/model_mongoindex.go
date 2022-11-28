package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoIndex struct {
	Key     *MongoIndexKeys    `json:"key"`
	Options *MongoIndexOptions `json:"options"`
}
