package usages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Usage struct {
	CurrentValue int64     `json:"currentValue"`
	Id           *string   `json:"id,omitempty"`
	Limit        int64     `json:"limit"`
	Name         UsageName `json:"name"`
	Unit         UsageUnit `json:"unit"`
}
