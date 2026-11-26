package nodepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Taint struct {
	Effect Effect  `json:"effect"`
	Key    string  `json:"key"`
	Value  *string `json:"value,omitempty"`
}
