package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridConnectionLimitsProperties struct {
	Current *int64 `json:"current,omitempty"`
	Maximum *int64 `json:"maximum,omitempty"`
}
