package tableservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableSignedIdentifier struct {
	AccessPolicy *TableAccessPolicy `json:"accessPolicy,omitempty"`
	Id           string             `json:"id"`
}
