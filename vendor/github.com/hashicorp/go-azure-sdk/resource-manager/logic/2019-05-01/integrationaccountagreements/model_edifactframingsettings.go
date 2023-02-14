package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdifactFramingSettings struct {
	CharacterEncoding               *string                 `json:"characterEncoding,omitempty"`
	CharacterSet                    EdifactCharacterSet     `json:"characterSet"`
	ComponentSeparator              int64                   `json:"componentSeparator"`
	DataElementSeparator            int64                   `json:"dataElementSeparator"`
	DecimalPointIndicator           EdifactDecimalIndicator `json:"decimalPointIndicator"`
	ProtocolVersion                 int64                   `json:"protocolVersion"`
	ReleaseIndicator                int64                   `json:"releaseIndicator"`
	RepetitionSeparator             int64                   `json:"repetitionSeparator"`
	SegmentTerminator               int64                   `json:"segmentTerminator"`
	SegmentTerminatorSuffix         SegmentTerminatorSuffix `json:"segmentTerminatorSuffix"`
	ServiceCodeListDirectoryVersion *string                 `json:"serviceCodeListDirectoryVersion,omitempty"`
}
