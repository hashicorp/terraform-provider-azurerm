package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12DelimiterOverrides struct {
	ComponentSeparator         int64                   `json:"componentSeparator"`
	DataElementSeparator       int64                   `json:"dataElementSeparator"`
	MessageId                  *string                 `json:"messageId,omitempty"`
	ProtocolVersion            *string                 `json:"protocolVersion,omitempty"`
	ReplaceCharacter           int64                   `json:"replaceCharacter"`
	ReplaceSeparatorsInPayload bool                    `json:"replaceSeparatorsInPayload"`
	SegmentTerminator          int64                   `json:"segmentTerminator"`
	SegmentTerminatorSuffix    SegmentTerminatorSuffix `json:"segmentTerminatorSuffix"`
	TargetNamespace            *string                 `json:"targetNamespace,omitempty"`
}
