package slice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Snssai struct {
	Sd  *string `json:"sd,omitempty"`
	Sst int64   `json:"sst"`
}
