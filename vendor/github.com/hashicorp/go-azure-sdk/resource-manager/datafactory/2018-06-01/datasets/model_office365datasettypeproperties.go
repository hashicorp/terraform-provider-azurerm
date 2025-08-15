package datasets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Office365DatasetTypeProperties struct {
	Predicate *interface{} `json:"predicate,omitempty"`
	TableName interface{}  `json:"tableName"`
}
