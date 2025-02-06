package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMachineLearningServiceOutputColumn struct {
	DataType *string `json:"dataType,omitempty"`
	MapTo    *int64  `json:"mapTo,omitempty"`
	Name     *string `json:"name,omitempty"`
}
