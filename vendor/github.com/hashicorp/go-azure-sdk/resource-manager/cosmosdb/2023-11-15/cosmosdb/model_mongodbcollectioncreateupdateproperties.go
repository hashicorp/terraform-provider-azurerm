package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBCollectionCreateUpdateProperties struct {
	Options  *CreateUpdateOptions      `json:"options,omitempty"`
	Resource MongoDBCollectionResource `json:"resource"`
}
