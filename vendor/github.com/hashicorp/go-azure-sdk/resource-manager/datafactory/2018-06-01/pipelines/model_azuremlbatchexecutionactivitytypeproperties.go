package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMLBatchExecutionActivityTypeProperties struct {
	GlobalParameters  *map[string]interface{}           `json:"globalParameters,omitempty"`
	WebServiceInputs  *map[string]AzureMLWebServiceFile `json:"webServiceInputs,omitempty"`
	WebServiceOutputs *map[string]AzureMLWebServiceFile `json:"webServiceOutputs,omitempty"`
}
