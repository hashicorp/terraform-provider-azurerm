package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMachineLearningServiceFunctionBindingProperties struct {
	ApiKey                   *string                                    `json:"apiKey,omitempty"`
	BatchSize                *int64                                     `json:"batchSize,omitempty"`
	Endpoint                 *string                                    `json:"endpoint,omitempty"`
	InputRequestName         *string                                    `json:"inputRequestName,omitempty"`
	Inputs                   *[]AzureMachineLearningServiceInputColumn  `json:"inputs,omitempty"`
	NumberOfParallelRequests *int64                                     `json:"numberOfParallelRequests,omitempty"`
	OutputResponseName       *string                                    `json:"outputResponseName,omitempty"`
	Outputs                  *[]AzureMachineLearningServiceOutputColumn `json:"outputs,omitempty"`
}
