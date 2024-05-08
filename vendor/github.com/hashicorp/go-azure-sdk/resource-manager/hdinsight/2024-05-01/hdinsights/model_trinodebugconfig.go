package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrinoDebugConfig struct {
	Enable  *bool  `json:"enable,omitempty"`
	Port    *int64 `json:"port,omitempty"`
	Suspend *bool  `json:"suspend,omitempty"`
}
