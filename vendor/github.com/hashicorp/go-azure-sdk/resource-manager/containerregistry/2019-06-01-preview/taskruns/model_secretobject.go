package taskruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretObject struct {
	Type  *SecretObjectType `json:"type,omitempty"`
	Value *string           `json:"value,omitempty"`
}
