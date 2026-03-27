package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SelfCheck struct {
	IntervalSeconds *int64           `json:"intervalSeconds,omitempty"`
	Mode            *OperationalMode `json:"mode,omitempty"`
	TimeoutSeconds  *int64           `json:"timeoutSeconds,omitempty"`
}
