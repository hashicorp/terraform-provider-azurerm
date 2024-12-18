package querykeys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryKey struct {
	Key  *string `json:"key,omitempty"`
	Name *string `json:"name,omitempty"`
}
