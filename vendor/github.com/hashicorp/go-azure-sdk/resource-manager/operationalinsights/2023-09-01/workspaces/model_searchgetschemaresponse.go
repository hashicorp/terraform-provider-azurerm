package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchGetSchemaResponse struct {
	Metadata *SearchMetadata      `json:"metadata,omitempty"`
	Value    *[]SearchSchemaValue `json:"value,omitempty"`
}
