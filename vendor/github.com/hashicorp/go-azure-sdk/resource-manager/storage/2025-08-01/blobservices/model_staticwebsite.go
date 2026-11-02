package blobservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticWebsite struct {
	DefaultIndexDocumentPath *string `json:"defaultIndexDocumentPath,omitempty"`
	Enabled                  bool    `json:"enabled"`
	ErrorDocument404Path     *string `json:"errorDocument404Path,omitempty"`
	IndexDocument            *string `json:"indexDocument,omitempty"`
}
