package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnStringInfo struct {
	ConnectionString *string               `json:"connectionString,omitempty"`
	Name             *string               `json:"name,omitempty"`
	Type             *ConnectionStringType `json:"type,omitempty"`
}
