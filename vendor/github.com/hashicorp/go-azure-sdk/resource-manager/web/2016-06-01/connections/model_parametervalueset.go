package connections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParameterValueSet struct {
	Name   *string                 `json:"name,omitempty"`
	Values *map[string]interface{} `json:"values,omitempty"`
}
