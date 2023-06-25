package integrationaccountschemas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentHash struct {
	Algorithm *string `json:"algorithm,omitempty"`
	Value     *string `json:"value,omitempty"`
}
