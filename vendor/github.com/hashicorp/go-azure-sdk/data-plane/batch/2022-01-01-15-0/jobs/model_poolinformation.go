package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolInformation struct {
	AutoPoolSpecification *AutoPoolSpecification `json:"autoPoolSpecification,omitempty"`
	PoolId                *string                `json:"poolId,omitempty"`
}
