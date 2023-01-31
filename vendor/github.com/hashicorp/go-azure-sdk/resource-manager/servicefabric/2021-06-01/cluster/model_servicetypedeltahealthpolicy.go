package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceTypeDeltaHealthPolicy struct {
	MaxPercentDeltaUnhealthyServices *int64 `json:"maxPercentDeltaUnhealthyServices,omitempty"`
}
