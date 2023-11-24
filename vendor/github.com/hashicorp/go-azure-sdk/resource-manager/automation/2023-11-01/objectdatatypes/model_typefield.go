package objectdatatypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TypeField struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}
