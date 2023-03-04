package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityCriteria struct {
	Type  *ConnectivityType `json:"type,omitempty"`
	Value *[]string         `json:"value,omitempty"`
}
