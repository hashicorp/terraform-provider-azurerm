package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDataExplorerOutputDataSourceProperties struct {
	AuthenticationMode *AuthenticationMode `json:"authenticationMode,omitempty"`
	Cluster            *string             `json:"cluster,omitempty"`
	Database           *string             `json:"database,omitempty"`
	Table              *string             `json:"table,omitempty"`
}
