package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConsistencyPolicy struct {
	DefaultConsistencyLevel DefaultConsistencyLevel `json:"defaultConsistencyLevel"`
	MaxIntervalInSeconds    *int64                  `json:"maxIntervalInSeconds,omitempty"`
	MaxStalenessPrefix      *int64                  `json:"maxStalenessPrefix,omitempty"`
}
