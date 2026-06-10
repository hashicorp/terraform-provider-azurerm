package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbCursorMethodsProperties struct {
	Limit   *interface{} `json:"limit,omitempty"`
	Project *interface{} `json:"project,omitempty"`
	Skip    *interface{} `json:"skip,omitempty"`
	Sort    *interface{} `json:"sort,omitempty"`
}
