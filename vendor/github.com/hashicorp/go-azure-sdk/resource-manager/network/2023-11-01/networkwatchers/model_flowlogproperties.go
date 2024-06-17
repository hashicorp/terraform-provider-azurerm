package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlowLogProperties struct {
	Enabled         bool                       `json:"enabled"`
	Format          *FlowLogFormatParameters   `json:"format,omitempty"`
	RetentionPolicy *RetentionPolicyParameters `json:"retentionPolicy,omitempty"`
	StorageId       string                     `json:"storageId"`
}
