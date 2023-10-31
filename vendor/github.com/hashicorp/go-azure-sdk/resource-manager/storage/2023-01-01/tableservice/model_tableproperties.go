package tableservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableProperties struct {
	SignedIdentifiers *[]TableSignedIdentifier `json:"signedIdentifiers,omitempty"`
	TableName         *string                  `json:"tableName,omitempty"`
}
