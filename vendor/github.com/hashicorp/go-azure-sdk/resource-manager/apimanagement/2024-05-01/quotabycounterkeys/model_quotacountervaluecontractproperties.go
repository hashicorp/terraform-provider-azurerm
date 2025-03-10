package quotabycounterkeys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaCounterValueContractProperties struct {
	CallsCount    *int64   `json:"callsCount,omitempty"`
	KbTransferred *float64 `json:"kbTransferred,omitempty"`
}
