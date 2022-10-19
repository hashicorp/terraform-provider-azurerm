package dataconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnectionCheckNameRequest struct {
	Name string             `json:"name"`
	Type DataConnectionType `json:"type"`
}
