package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12ProtocolSettings struct {
	AcknowledgementSettings X12AcknowledgementSettings `json:"acknowledgementSettings"`
	EnvelopeOverrides       *[]X12EnvelopeOverride     `json:"envelopeOverrides,omitempty"`
	EnvelopeSettings        X12EnvelopeSettings        `json:"envelopeSettings"`
	FramingSettings         X12FramingSettings         `json:"framingSettings"`
	MessageFilter           X12MessageFilter           `json:"messageFilter"`
	MessageFilterList       *[]X12MessageIdentifier    `json:"messageFilterList,omitempty"`
	ProcessingSettings      X12ProcessingSettings      `json:"processingSettings"`
	SchemaReferences        []X12SchemaReference       `json:"schemaReferences"`
	SecuritySettings        X12SecuritySettings        `json:"securitySettings"`
	ValidationOverrides     *[]X12ValidationOverride   `json:"validationOverrides,omitempty"`
	ValidationSettings      X12ValidationSettings      `json:"validationSettings"`
	X12DelimiterOverrides   *[]X12DelimiterOverrides   `json:"x12DelimiterOverrides,omitempty"`
}
