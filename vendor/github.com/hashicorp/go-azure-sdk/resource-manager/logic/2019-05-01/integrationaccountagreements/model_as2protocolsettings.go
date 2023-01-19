package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AS2ProtocolSettings struct {
	AcknowledgementConnectionSettings AS2AcknowledgementConnectionSettings `json:"acknowledgementConnectionSettings"`
	EnvelopeSettings                  AS2EnvelopeSettings                  `json:"envelopeSettings"`
	ErrorSettings                     AS2ErrorSettings                     `json:"errorSettings"`
	MdnSettings                       AS2MdnSettings                       `json:"mdnSettings"`
	MessageConnectionSettings         AS2MessageConnectionSettings         `json:"messageConnectionSettings"`
	SecuritySettings                  AS2SecuritySettings                  `json:"securitySettings"`
	ValidationSettings                AS2ValidationSettings                `json:"validationSettings"`
}
