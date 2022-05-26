package extensionoperationstatus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationStatusResult struct {
	Error      *ErrorDetail       `json:"error,omitempty"`
	Id         *string            `json:"id,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties *map[string]string `json:"properties,omitempty"`
	Status     string             `json:"status"`
}
