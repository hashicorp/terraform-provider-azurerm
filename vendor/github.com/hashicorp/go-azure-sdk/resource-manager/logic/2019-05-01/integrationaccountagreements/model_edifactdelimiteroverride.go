package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdifactDelimiterOverride struct {
	ComponentSeparator             int64                   `json:"componentSeparator"`
	DataElementSeparator           int64                   `json:"dataElementSeparator"`
	DecimalPointIndicator          EdifactDecimalIndicator `json:"decimalPointIndicator"`
	MessageAssociationAssignedCode *string                 `json:"messageAssociationAssignedCode,omitempty"`
	MessageId                      *string                 `json:"messageId,omitempty"`
	MessageRelease                 *string                 `json:"messageRelease,omitempty"`
	MessageVersion                 *string                 `json:"messageVersion,omitempty"`
	ReleaseIndicator               int64                   `json:"releaseIndicator"`
	RepetitionSeparator            int64                   `json:"repetitionSeparator"`
	SegmentTerminator              int64                   `json:"segmentTerminator"`
	SegmentTerminatorSuffix        SegmentTerminatorSuffix `json:"segmentTerminatorSuffix"`
	TargetNamespace                *string                 `json:"targetNamespace,omitempty"`
}
