package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SshProfile struct {
	Count     int64   `json:"count"`
	PodPrefix *string `json:"podPrefix,omitempty"`
	VMSize    *string `json:"vmSize,omitempty"`
}
