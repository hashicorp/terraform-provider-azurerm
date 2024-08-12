package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseAccountConnectionString struct {
	ConnectionString *string `json:"connectionString,omitempty"`
	Description      *string `json:"description,omitempty"`
	KeyKind          *Kind   `json:"keyKind,omitempty"`
	Type             *Type   `json:"type,omitempty"`
}
