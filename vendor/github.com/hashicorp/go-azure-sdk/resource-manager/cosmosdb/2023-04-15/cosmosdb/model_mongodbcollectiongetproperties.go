package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBCollectionGetProperties struct {
	Options  *OptionsResource                        `json:"options,omitempty"`
	Resource *MongoDBCollectionGetPropertiesResource `json:"resource,omitempty"`
}
