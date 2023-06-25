package functions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMachineLearningWebServiceFunctionBindingProperties struct {
	ApiKey    *string                                       `json:"apiKey,omitempty"`
	BatchSize *int64                                        `json:"batchSize,omitempty"`
	Endpoint  *string                                       `json:"endpoint,omitempty"`
	Inputs    *AzureMachineLearningWebServiceInputs         `json:"inputs,omitempty"`
	Outputs   *[]AzureMachineLearningWebServiceOutputColumn `json:"outputs,omitempty"`
}
