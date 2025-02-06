package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DocumentDbOutputDataSourceProperties struct {
	AccountId             *string             `json:"accountId,omitempty"`
	AccountKey            *string             `json:"accountKey,omitempty"`
	AuthenticationMode    *AuthenticationMode `json:"authenticationMode,omitempty"`
	CollectionNamePattern *string             `json:"collectionNamePattern,omitempty"`
	Database              *string             `json:"database,omitempty"`
	DocumentId            *string             `json:"documentId,omitempty"`
	PartitionKey          *string             `json:"partitionKey,omitempty"`
}
