package elasticpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticPoolPerDatabaseSettings struct {
	AutoPauseDelay *int64   `json:"autoPauseDelay,omitempty"`
	MaxCapacity    *float64 `json:"maxCapacity,omitempty"`
	MinCapacity    *float64 `json:"minCapacity,omitempty"`
}
