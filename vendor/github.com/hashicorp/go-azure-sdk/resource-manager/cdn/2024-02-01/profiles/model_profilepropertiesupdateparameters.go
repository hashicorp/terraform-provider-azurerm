package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProfilePropertiesUpdateParameters struct {
	LogScrubbing                 *ProfileLogScrubbing `json:"logScrubbing,omitempty"`
	OriginResponseTimeoutSeconds *int64               `json:"originResponseTimeoutSeconds,omitempty"`
}
