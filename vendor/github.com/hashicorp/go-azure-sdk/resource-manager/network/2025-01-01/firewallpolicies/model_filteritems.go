package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FilterItems struct {
	Field  *string   `json:"field,omitempty"`
	Values *[]string `json:"values,omitempty"`
}
