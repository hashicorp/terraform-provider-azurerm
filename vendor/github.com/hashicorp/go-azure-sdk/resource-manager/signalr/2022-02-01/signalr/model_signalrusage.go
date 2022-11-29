package signalr

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SignalRUsage struct {
	CurrentValue *int64            `json:"currentValue,omitempty"`
	Id           *string           `json:"id,omitempty"`
	Limit        *int64            `json:"limit,omitempty"`
	Name         *SignalRUsageName `json:"name,omitempty"`
	Unit         *string           `json:"unit,omitempty"`
}
