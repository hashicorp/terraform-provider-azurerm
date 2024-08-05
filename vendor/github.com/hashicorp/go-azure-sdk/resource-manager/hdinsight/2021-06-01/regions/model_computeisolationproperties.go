package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeIsolationProperties struct {
	EnableComputeIsolation *bool   `json:"enableComputeIsolation,omitempty"`
	HostSku                *string `json:"hostSku,omitempty"`
}
