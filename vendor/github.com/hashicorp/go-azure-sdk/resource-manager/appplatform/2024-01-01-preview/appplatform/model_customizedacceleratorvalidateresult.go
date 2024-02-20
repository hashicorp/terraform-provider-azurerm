package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomizedAcceleratorValidateResult struct {
	ErrorMessage *string                                   `json:"errorMessage,omitempty"`
	State        *CustomizedAcceleratorValidateResultState `json:"state,omitempty"`
}
