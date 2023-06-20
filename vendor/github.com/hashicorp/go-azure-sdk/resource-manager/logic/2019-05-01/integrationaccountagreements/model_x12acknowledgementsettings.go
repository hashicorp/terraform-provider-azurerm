package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12AcknowledgementSettings struct {
	AcknowledgementControlNumberLowerBound int64   `json:"acknowledgementControlNumberLowerBound"`
	AcknowledgementControlNumberPrefix     *string `json:"acknowledgementControlNumberPrefix,omitempty"`
	AcknowledgementControlNumberSuffix     *string `json:"acknowledgementControlNumberSuffix,omitempty"`
	AcknowledgementControlNumberUpperBound int64   `json:"acknowledgementControlNumberUpperBound"`
	BatchFunctionalAcknowledgements        bool    `json:"batchFunctionalAcknowledgements"`
	BatchImplementationAcknowledgements    bool    `json:"batchImplementationAcknowledgements"`
	BatchTechnicalAcknowledgements         bool    `json:"batchTechnicalAcknowledgements"`
	FunctionalAcknowledgementVersion       *string `json:"functionalAcknowledgementVersion,omitempty"`
	ImplementationAcknowledgementVersion   *string `json:"implementationAcknowledgementVersion,omitempty"`
	NeedFunctionalAcknowledgement          bool    `json:"needFunctionalAcknowledgement"`
	NeedImplementationAcknowledgement      bool    `json:"needImplementationAcknowledgement"`
	NeedLoopForValidMessages               bool    `json:"needLoopForValidMessages"`
	NeedTechnicalAcknowledgement           bool    `json:"needTechnicalAcknowledgement"`
	RolloverAcknowledgementControlNumber   bool    `json:"rolloverAcknowledgementControlNumber"`
	SendSynchronousAcknowledgement         bool    `json:"sendSynchronousAcknowledgement"`
}
