package recordsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SrvRecord struct {
	Port     *int64  `json:"port,omitempty"`
	Priority *int64  `json:"priority,omitempty"`
	Target   *string `json:"target,omitempty"`
	Weight   *int64  `json:"weight,omitempty"`
}
