package threatintelligence

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceParsedPattern struct {
	PatternTypeKey    *string                                     `json:"patternTypeKey,omitempty"`
	PatternTypeValues *[]ThreatIntelligenceParsedPatternTypeValue `json:"patternTypeValues,omitempty"`
}
