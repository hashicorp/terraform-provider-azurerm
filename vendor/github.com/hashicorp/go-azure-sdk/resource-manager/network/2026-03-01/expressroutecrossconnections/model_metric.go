package expressroutecrossconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Metric struct {
	Name  *string  `json:"name,omitempty"`
	Unit  *string  `json:"unit,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
