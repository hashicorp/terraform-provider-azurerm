package scenvironmentrecords

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentProperties struct {
	Metadata               *SCMetadataEntity       `json:"metadata,omitempty"`
	StreamGovernanceConfig *StreamGovernanceConfig `json:"streamGovernanceConfig,omitempty"`
}
