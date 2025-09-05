package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraError struct {
	AdditionalErrorInfo *string `json:"additionalErrorInfo,omitempty"`
	Code                *string `json:"code,omitempty"`
	Message             *string `json:"message,omitempty"`
	Target              *string `json:"target,omitempty"`
}
