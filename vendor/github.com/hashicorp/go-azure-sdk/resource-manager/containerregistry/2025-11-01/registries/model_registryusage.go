package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryUsage struct {
	CurrentValue *int64             `json:"currentValue,omitempty"`
	Limit        *int64             `json:"limit,omitempty"`
	Name         *string            `json:"name,omitempty"`
	Unit         *RegistryUsageUnit `json:"unit,omitempty"`
}
