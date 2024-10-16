package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbError struct {
	Code    *string           `json:"code,omitempty"`
	Count   *int64            `json:"count,omitempty"`
	Message *string           `json:"message,omitempty"`
	Type    *MongoDbErrorType `json:"type,omitempty"`
}
