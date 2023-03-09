package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFunctionOutputDataSourceProperties struct {
	ApiKey          *string  `json:"apiKey,omitempty"`
	FunctionAppName *string  `json:"functionAppName,omitempty"`
	FunctionName    *string  `json:"functionName,omitempty"`
	MaxBatchCount   *float64 `json:"maxBatchCount,omitempty"`
	MaxBatchSize    *float64 `json:"maxBatchSize,omitempty"`
}
