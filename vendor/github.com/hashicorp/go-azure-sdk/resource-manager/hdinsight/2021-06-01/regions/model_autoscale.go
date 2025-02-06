package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Autoscale struct {
	Capacity   *AutoscaleCapacity   `json:"capacity,omitempty"`
	Recurrence *AutoscaleRecurrence `json:"recurrence,omitempty"`
}
