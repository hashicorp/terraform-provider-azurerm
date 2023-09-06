package databases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Persistence struct {
	AofEnabled   *bool         `json:"aofEnabled,omitempty"`
	AofFrequency *AofFrequency `json:"aofFrequency,omitempty"`
	RdbEnabled   *bool         `json:"rdbEnabled,omitempty"`
	RdbFrequency *RdbFrequency `json:"rdbFrequency,omitempty"`
}
