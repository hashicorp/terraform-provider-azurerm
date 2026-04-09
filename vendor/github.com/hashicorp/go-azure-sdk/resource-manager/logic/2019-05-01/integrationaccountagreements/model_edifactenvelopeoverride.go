package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdifactEnvelopeOverride struct {
	ApplicationPassword            *string `json:"applicationPassword,omitempty"`
	AssociationAssignedCode        *string `json:"associationAssignedCode,omitempty"`
	ControllingAgencyCode          *string `json:"controllingAgencyCode,omitempty"`
	FunctionalGroupId              *string `json:"functionalGroupId,omitempty"`
	GroupHeaderMessageRelease      *string `json:"groupHeaderMessageRelease,omitempty"`
	GroupHeaderMessageVersion      *string `json:"groupHeaderMessageVersion,omitempty"`
	MessageAssociationAssignedCode *string `json:"messageAssociationAssignedCode,omitempty"`
	MessageId                      *string `json:"messageId,omitempty"`
	MessageRelease                 *string `json:"messageRelease,omitempty"`
	MessageVersion                 *string `json:"messageVersion,omitempty"`
	ReceiverApplicationId          *string `json:"receiverApplicationId,omitempty"`
	ReceiverApplicationQualifier   *string `json:"receiverApplicationQualifier,omitempty"`
	SenderApplicationId            *string `json:"senderApplicationId,omitempty"`
	SenderApplicationQualifier     *string `json:"senderApplicationQualifier,omitempty"`
	TargetNamespace                *string `json:"targetNamespace,omitempty"`
}
