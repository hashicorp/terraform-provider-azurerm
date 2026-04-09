package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12FramingSettings struct {
	CharacterSet               X12CharacterSet         `json:"characterSet"`
	ComponentSeparator         int64                   `json:"componentSeparator"`
	DataElementSeparator       int64                   `json:"dataElementSeparator"`
	ReplaceCharacter           int64                   `json:"replaceCharacter"`
	ReplaceSeparatorsInPayload bool                    `json:"replaceSeparatorsInPayload"`
	SegmentTerminator          int64                   `json:"segmentTerminator"`
	SegmentTerminatorSuffix    SegmentTerminatorSuffix `json:"segmentTerminatorSuffix"`
}
