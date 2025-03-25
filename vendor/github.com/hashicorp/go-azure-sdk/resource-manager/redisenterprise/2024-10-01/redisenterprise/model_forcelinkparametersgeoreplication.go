package redisenterprise

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForceLinkParametersGeoReplication struct {
	GroupNickname   *string           `json:"groupNickname,omitempty"`
	LinkedDatabases *[]LinkedDatabase `json:"linkedDatabases,omitempty"`
}
