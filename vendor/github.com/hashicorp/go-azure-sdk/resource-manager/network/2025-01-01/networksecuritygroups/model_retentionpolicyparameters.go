package networksecuritygroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionPolicyParameters struct {
	Days    *int64 `json:"days,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}
