package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostsGetRequest struct {
	UserEmail string    `json:"userEmail"`
	VMIds     *[]string `json:"vmIds,omitempty"`
}
