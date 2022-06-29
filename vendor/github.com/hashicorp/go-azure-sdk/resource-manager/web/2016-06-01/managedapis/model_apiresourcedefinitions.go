package managedapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiResourceDefinitions struct {
	ModifiedSwaggerUrl *string `json:"modifiedSwaggerUrl,omitempty"`
	OriginalSwaggerUrl *string `json:"originalSwaggerUrl,omitempty"`
}
