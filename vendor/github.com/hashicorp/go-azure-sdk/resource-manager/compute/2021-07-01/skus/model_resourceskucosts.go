package skus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceSkuCosts struct {
	ExtendedUnit *string `json:"extendedUnit,omitempty"`
	MeterID      *string `json:"meterID,omitempty"`
	Quantity     *int64  `json:"quantity,omitempty"`
}
