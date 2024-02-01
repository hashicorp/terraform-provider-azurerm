package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMachineLearningStudioInputs struct {
	ColumnNames *[]AzureMachineLearningStudioInputColumn `json:"columnNames,omitempty"`
	Name        *string                                  `json:"name,omitempty"`
}
