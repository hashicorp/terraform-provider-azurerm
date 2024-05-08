package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterInstanceViewStatus struct {
	Message *string `json:"message,omitempty"`
	Ready   string  `json:"ready"`
	Reason  *string `json:"reason,omitempty"`
}
