package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdifactEnvelopeSettings struct {
	ApplicationReferenceId                       *string `json:"applicationReferenceId,omitempty"`
	ApplyDelimiterStringAdvice                   bool    `json:"applyDelimiterStringAdvice"`
	CommunicationAgreementId                     *string `json:"communicationAgreementId,omitempty"`
	CreateGroupingSegments                       bool    `json:"createGroupingSegments"`
	EnableDefaultGroupHeaders                    bool    `json:"enableDefaultGroupHeaders"`
	FunctionalGroupId                            *string `json:"functionalGroupId,omitempty"`
	GroupApplicationPassword                     *string `json:"groupApplicationPassword,omitempty"`
	GroupApplicationReceiverId                   *string `json:"groupApplicationReceiverId,omitempty"`
	GroupApplicationReceiverQualifier            *string `json:"groupApplicationReceiverQualifier,omitempty"`
	GroupApplicationSenderId                     *string `json:"groupApplicationSenderId,omitempty"`
	GroupApplicationSenderQualifier              *string `json:"groupApplicationSenderQualifier,omitempty"`
	GroupAssociationAssignedCode                 *string `json:"groupAssociationAssignedCode,omitempty"`
	GroupControlNumberLowerBound                 int64   `json:"groupControlNumberLowerBound"`
	GroupControlNumberPrefix                     *string `json:"groupControlNumberPrefix,omitempty"`
	GroupControlNumberSuffix                     *string `json:"groupControlNumberSuffix,omitempty"`
	GroupControlNumberUpperBound                 int64   `json:"groupControlNumberUpperBound"`
	GroupControllingAgencyCode                   *string `json:"groupControllingAgencyCode,omitempty"`
	GroupMessageRelease                          *string `json:"groupMessageRelease,omitempty"`
	GroupMessageVersion                          *string `json:"groupMessageVersion,omitempty"`
	InterchangeControlNumberLowerBound           int64   `json:"interchangeControlNumberLowerBound"`
	InterchangeControlNumberPrefix               *string `json:"interchangeControlNumberPrefix,omitempty"`
	InterchangeControlNumberSuffix               *string `json:"interchangeControlNumberSuffix,omitempty"`
	InterchangeControlNumberUpperBound           int64   `json:"interchangeControlNumberUpperBound"`
	IsTestInterchange                            bool    `json:"isTestInterchange"`
	OverwriteExistingTransactionSetControlNumber bool    `json:"overwriteExistingTransactionSetControlNumber"`
	ProcessingPriorityCode                       *string `json:"processingPriorityCode,omitempty"`
	ReceiverInternalIdentification               *string `json:"receiverInternalIdentification,omitempty"`
	ReceiverInternalSubIdentification            *string `json:"receiverInternalSubIdentification,omitempty"`
	ReceiverReverseRoutingAddress                *string `json:"receiverReverseRoutingAddress,omitempty"`
	RecipientReferencePasswordQualifier          *string `json:"recipientReferencePasswordQualifier,omitempty"`
	RecipientReferencePasswordValue              *string `json:"recipientReferencePasswordValue,omitempty"`
	RolloverGroupControlNumber                   bool    `json:"rolloverGroupControlNumber"`
	RolloverInterchangeControlNumber             bool    `json:"rolloverInterchangeControlNumber"`
	RolloverTransactionSetControlNumber          bool    `json:"rolloverTransactionSetControlNumber"`
	SenderInternalIdentification                 *string `json:"senderInternalIdentification,omitempty"`
	SenderInternalSubIdentification              *string `json:"senderInternalSubIdentification,omitempty"`
	SenderReverseRoutingAddress                  *string `json:"senderReverseRoutingAddress,omitempty"`
	TransactionSetControlNumberLowerBound        int64   `json:"transactionSetControlNumberLowerBound"`
	TransactionSetControlNumberPrefix            *string `json:"transactionSetControlNumberPrefix,omitempty"`
	TransactionSetControlNumberSuffix            *string `json:"transactionSetControlNumberSuffix,omitempty"`
	TransactionSetControlNumberUpperBound        int64   `json:"transactionSetControlNumberUpperBound"`
}
