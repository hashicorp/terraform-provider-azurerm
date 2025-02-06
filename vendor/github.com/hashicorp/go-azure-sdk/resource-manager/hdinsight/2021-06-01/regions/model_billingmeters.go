package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingMeters struct {
	Meter          *string `json:"meter,omitempty"`
	MeterParameter *string `json:"meterParameter,omitempty"`
	Unit           *string `json:"unit,omitempty"`
}
