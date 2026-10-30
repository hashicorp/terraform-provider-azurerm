package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceUsage struct {
	CurrentValue *int64             `json:"currentValue,omitempty"`
	Limit        *int64             `json:"limit,omitempty"`
	ResourceType *string            `json:"resourceType,omitempty"`
	Unit         *ResourceUsageUnit `json:"unit,omitempty"`
}
