package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThroughputPolicyResource struct {
	IncrementPercent *int64 `json:"incrementPercent,omitempty"`
	IsEnabled        *bool  `json:"isEnabled,omitempty"`
}
