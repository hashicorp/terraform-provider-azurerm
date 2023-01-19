package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12EnvelopeOverride struct {
	DateFormat               X12DateFormat `json:"dateFormat"`
	FunctionalIdentifierCode *string       `json:"functionalIdentifierCode,omitempty"`
	HeaderVersion            string        `json:"headerVersion"`
	MessageId                string        `json:"messageId"`
	ProtocolVersion          string        `json:"protocolVersion"`
	ReceiverApplicationId    string        `json:"receiverApplicationId"`
	ResponsibleAgencyCode    string        `json:"responsibleAgencyCode"`
	SenderApplicationId      string        `json:"senderApplicationId"`
	TargetNamespace          string        `json:"targetNamespace"`
	TimeFormat               X12TimeFormat `json:"timeFormat"`
}
