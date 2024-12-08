package recordsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Digest struct {
	AlgorithmType *int64  `json:"algorithmType,omitempty"`
	Value         *string `json:"value,omitempty"`
}
