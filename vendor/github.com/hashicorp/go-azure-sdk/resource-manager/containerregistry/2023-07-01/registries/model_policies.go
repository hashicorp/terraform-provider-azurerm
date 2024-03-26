package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Policies struct {
	ExportPolicy     *ExportPolicy     `json:"exportPolicy,omitempty"`
	QuarantinePolicy *QuarantinePolicy `json:"quarantinePolicy,omitempty"`
	RetentionPolicy  *RetentionPolicy  `json:"retentionPolicy,omitempty"`
	TrustPolicy      *TrustPolicy      `json:"trustPolicy,omitempty"`
}
