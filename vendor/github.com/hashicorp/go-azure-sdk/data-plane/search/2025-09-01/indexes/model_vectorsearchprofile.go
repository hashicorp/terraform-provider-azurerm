package indexes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VectorSearchProfile struct {
	Algorithm   string  `json:"algorithm"`
	Compression *string `json:"compression,omitempty"`
	Name        string  `json:"name"`
	Vectorizer  *string `json:"vectorizer,omitempty"`
}
