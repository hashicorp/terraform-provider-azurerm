package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12EnvelopeSettings struct {
	ControlStandardsId                           int64          `json:"controlStandardsId"`
	ControlVersionNumber                         string         `json:"controlVersionNumber"`
	EnableDefaultGroupHeaders                    bool           `json:"enableDefaultGroupHeaders"`
	FunctionalGroupId                            *string        `json:"functionalGroupId,omitempty"`
	GroupControlNumberLowerBound                 int64          `json:"groupControlNumberLowerBound"`
	GroupControlNumberUpperBound                 int64          `json:"groupControlNumberUpperBound"`
	GroupHeaderAgencyCode                        string         `json:"groupHeaderAgencyCode"`
	GroupHeaderDateFormat                        X12DateFormat  `json:"groupHeaderDateFormat"`
	GroupHeaderTimeFormat                        X12TimeFormat  `json:"groupHeaderTimeFormat"`
	GroupHeaderVersion                           string         `json:"groupHeaderVersion"`
	InterchangeControlNumberLowerBound           int64          `json:"interchangeControlNumberLowerBound"`
	InterchangeControlNumberUpperBound           int64          `json:"interchangeControlNumberUpperBound"`
	OverwriteExistingTransactionSetControlNumber bool           `json:"overwriteExistingTransactionSetControlNumber"`
	ReceiverApplicationId                        string         `json:"receiverApplicationId"`
	RolloverGroupControlNumber                   bool           `json:"rolloverGroupControlNumber"`
	RolloverInterchangeControlNumber             bool           `json:"rolloverInterchangeControlNumber"`
	RolloverTransactionSetControlNumber          bool           `json:"rolloverTransactionSetControlNumber"`
	SenderApplicationId                          string         `json:"senderApplicationId"`
	TransactionSetControlNumberLowerBound        int64          `json:"transactionSetControlNumberLowerBound"`
	TransactionSetControlNumberPrefix            *string        `json:"transactionSetControlNumberPrefix,omitempty"`
	TransactionSetControlNumberSuffix            *string        `json:"transactionSetControlNumberSuffix,omitempty"`
	TransactionSetControlNumberUpperBound        int64          `json:"transactionSetControlNumberUpperBound"`
	UsageIndicator                               UsageIndicator `json:"usageIndicator"`
	UseControlStandardsIdAsRepetitionCharacter   bool           `json:"useControlStandardsIdAsRepetitionCharacter"`
}
