package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdifactProtocolSettings struct {
	AcknowledgementSettings   EdifactAcknowledgementSettings `json:"acknowledgementSettings"`
	EdifactDelimiterOverrides *[]EdifactDelimiterOverride    `json:"edifactDelimiterOverrides,omitempty"`
	EnvelopeOverrides         *[]EdifactEnvelopeOverride     `json:"envelopeOverrides,omitempty"`
	EnvelopeSettings          EdifactEnvelopeSettings        `json:"envelopeSettings"`
	FramingSettings           EdifactFramingSettings         `json:"framingSettings"`
	MessageFilter             EdifactMessageFilter           `json:"messageFilter"`
	MessageFilterList         *[]EdifactMessageIdentifier    `json:"messageFilterList,omitempty"`
	ProcessingSettings        EdifactProcessingSettings      `json:"processingSettings"`
	SchemaReferences          []EdifactSchemaReference       `json:"schemaReferences"`
	ValidationOverrides       *[]EdifactValidationOverride   `json:"validationOverrides,omitempty"`
	ValidationSettings        EdifactValidationSettings      `json:"validationSettings"`
}
