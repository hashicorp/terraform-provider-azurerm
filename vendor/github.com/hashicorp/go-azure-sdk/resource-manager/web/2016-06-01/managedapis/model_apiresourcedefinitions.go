package managedapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiResourceDefinitions struct {
	ModifiedSwaggerURL *string `json:"modifiedSwaggerUrl,omitempty"`
	OriginalSwaggerURL *string `json:"originalSwaggerUrl,omitempty"`
}
