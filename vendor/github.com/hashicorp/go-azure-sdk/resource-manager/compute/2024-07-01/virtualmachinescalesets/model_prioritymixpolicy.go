package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PriorityMixPolicy struct {
	BaseRegularPriorityCount           *int64 `json:"baseRegularPriorityCount,omitempty"`
	RegularPriorityPercentageAboveBase *int64 `json:"regularPriorityPercentageAboveBase,omitempty"`
}
