package assignment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParameterValue struct {
	Reference *SecretValueReference `json:"reference,omitempty"`
	Value     *interface{}          `json:"value,omitempty"`
}
