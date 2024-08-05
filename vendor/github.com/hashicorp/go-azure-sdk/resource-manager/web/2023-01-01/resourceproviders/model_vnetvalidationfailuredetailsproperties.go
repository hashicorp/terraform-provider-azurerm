package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VnetValidationFailureDetailsProperties struct {
	Failed      *bool                        `json:"failed,omitempty"`
	FailedTests *[]VnetValidationTestFailure `json:"failedTests,omitempty"`
	Message     *string                      `json:"message,omitempty"`
	Warnings    *[]VnetValidationTestFailure `json:"warnings,omitempty"`
}
