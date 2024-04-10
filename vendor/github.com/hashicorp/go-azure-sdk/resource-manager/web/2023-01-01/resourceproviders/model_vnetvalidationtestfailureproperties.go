package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VnetValidationTestFailureProperties struct {
	Details  *string `json:"details,omitempty"`
	TestName *string `json:"testName,omitempty"`
}
