package catalogs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GitCatalog struct {
	Branch           *string `json:"branch,omitempty"`
	Path             *string `json:"path,omitempty"`
	SecretIdentifier *string `json:"secretIdentifier,omitempty"`
	Uri              *string `json:"uri,omitempty"`
}
