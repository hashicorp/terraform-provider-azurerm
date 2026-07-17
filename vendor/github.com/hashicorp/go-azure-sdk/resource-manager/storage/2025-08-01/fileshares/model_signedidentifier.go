package fileshares

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SignedIdentifier struct {
	AccessPolicy *AccessPolicy `json:"accessPolicy,omitempty"`
	Id           *string       `json:"id,omitempty"`
}
