package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMachineLearningStudioFunctionBindingProperties struct {
	ApiKey    *string                                   `json:"apiKey,omitempty"`
	BatchSize *int64                                    `json:"batchSize,omitempty"`
	Endpoint  *string                                   `json:"endpoint,omitempty"`
	Inputs    *AzureMachineLearningStudioInputs         `json:"inputs,omitempty"`
	Outputs   *[]AzureMachineLearningStudioOutputColumn `json:"outputs,omitempty"`
}
