package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureSqlReferenceInputDataSourceProperties struct {
	AuthenticationMode *AuthenticationMode `json:"authenticationMode,omitempty"`
	Database           *string             `json:"database,omitempty"`
	DeltaSnapshotQuery *string             `json:"deltaSnapshotQuery,omitempty"`
	FullSnapshotQuery  *string             `json:"fullSnapshotQuery,omitempty"`
	Password           *string             `json:"password,omitempty"`
	RefreshRate        *string             `json:"refreshRate,omitempty"`
	RefreshType        *RefreshType        `json:"refreshType,omitempty"`
	Server             *string             `json:"server,omitempty"`
	User               *string             `json:"user,omitempty"`
}
