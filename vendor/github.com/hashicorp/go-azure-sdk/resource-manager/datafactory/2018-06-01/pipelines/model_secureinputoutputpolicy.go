package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecureInputOutputPolicy struct {
	SecureInput  *bool `json:"secureInput,omitempty"`
	SecureOutput *bool `json:"secureOutput,omitempty"`
}
